import sqlite3
from sqlalchemy import create_engine

engine = create_engine("sqlite:///f1.db", echo=True)

import os
from dotenv import load_dotenv

load_dotenv()

# from langchain_openai import ChatOpenAI #OpenAI LLM
from langchain_google_genai import ChatGoogleGenerativeAI  # Google LLM


# -------- Chain de LLM para validar y reformular prompts ---------
# from langchain.chains import LLMChain
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain

prompt_template = """Eres un experto en SQL.
Tu tarea es verificar la validez de un prompt dado por un usuario y, si es necesario, reformularlo para que sea más claro y específico.

- Si la consulta original es ambigua o puede interpretarse de más de una forma, hacé preguntas aclaratorias.
- Si la consulta es clara y no necesita aclaración, devolvé exactamente: NO_CLARIFICATION_NEEDED

Ejemplo:
Usuario: ¿Cuántos atendió Juan?
Asistente:
1. ¿A qué se refiere con "atendió"? (consultas, estudios, turnos, etc.)
2. ¿Quién es "Juan"? ¿Tenés apellido o rol (médico, paciente)?
3. ¿Querés filtrar por fechas?

Ejemplo de pregunta clara:
Usuario: ¿Cuántos pacientes atendió Juan Pérez en 2023?
Asistente: NO_CLARIFICATION_NEEDED

Generá preguntas cortas y claras para que el usuario aclare su intención, una por línea. No respondas la consulta.

Usuario: la siguiente consulta puede ser ambigua: "{pregunta}"
Asistente:
"""
clarificador_prompt = PromptTemplate(
    input_variables=["pregunta"],
    template=prompt_template,
)
# llm = ChatOpenAI(model="gpt-4o", temperature=0, openai_api_key=os.getenv("OPENAI_API_KEY")) #openai llm
llm = ChatGoogleGenerativeAI(
    model="gemini-2.5-flash",
    temperature=0.2,
    google_api_key=os.getenv("GOOGLE_API_KEY"),
)  # google llm

clarificador_chain = LLMChain(
    llm=llm,
    prompt=clarificador_prompt,
    verbose=True,
)

from langchain_community.utilities import SQLDatabase
from langchain_community.agent_toolkits import create_sql_agent
from langchain.agents import AgentType

db = SQLDatabase(engine=engine)
# --------- Agente SQL con LLM ---------
# llm = ChatOpenAI(model="gpt-4o", temperature=0, openai_api_key=os.getenv("OPENAI_API_KEY")) #openai llm
# agente = create_sql_agent(llm=llm, database=db,agent_type="openai-tools" , verbose=True) #verbose=True para ver como "piensa" el agente

llm = ChatGoogleGenerativeAI(
    model="gemini-2.5-flash", temperature=0, google_api_key=os.getenv("GOOGLE_API_KEY")
)  # google llm
# Al parecer gemini-2.5-pro es mas lenta que gemini-2.0-flash y mucho mas lenta que 2.5-flash
agente = create_sql_agent(
    llm=llm,
    db=db,
    agent_type=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True,
    handle_parsing_errors=True,
)

# -------- LLM Chain para explicar la consulta SQL generada ---------
# 6. Chain explicador
template_explicador = PromptTemplate.from_template(
    """
Tenés que explicarle al usuario un resultado de una consulta SQL que pidió en lenguaje natural.

Pregunta original:
"{pregunta}"

Aclaraciones:
{aclaraciones}

Resultado de la consulta SQL:
"{resultado}"

Respondé con una frase como:
"La respuesta es: ..." y luego explicá en lenguaje claro el significado de ese resultado, como si se lo explicaras a alguien sin conocimientos técnicos.
"""
)

explicador_chain = LLMChain(llm=llm, prompt=template_explicador)


template = """
Sos un asistente que clasifica si una explicación fue útil para el usuario.

Respuesta del usuario:
"{respuesta_usuario}"

Clasificá esta respuesta como una de las siguientes opciones (solo una palabra):
- útil
- no útil
"""

clasificador_prompt = PromptTemplate.from_template(template)
clasificador_chain = LLMChain(prompt=clasificador_prompt, llm=llm)

template = """
Tenés una conversación previa con el usuario, en la que se intentó responder una pregunta en lenguaje natural transformándola en SQL. A continuación se incluye el historial y un comentario final del usuario.

Historial:
{historial}

Nueva aclaración o corrección del usuario:
{nueva_aclaracion}

Pregunta original:
{pregunta_original}

Reformulá una nueva pregunta clara, específica y completa en lenguaje natural que tenga en cuenta todo el contexto y la aclaración.
Solo devolvé la nueva pregunta, sin explicaciones adicionales.
"""

reformulador_prompt = PromptTemplate.from_template(template)
reformulador_chain = LLMChain(prompt=reformulador_prompt, llm=llm)

correction_prompt = PromptTemplate(
    input_variables=["query", "error"],
    template="""
El siguiente query en lenguaje natural produjo un error al ser ejecutado por un agente SQL:

Query:
{query}

Error:
{error}

Corrige el query para que funcione correctamente y respete el esquema. Devuelve solo el nuevo query.
""",
)

correction_chain = LLMChain(llm=llm, prompt=correction_prompt)


def loop_consulta_sql(
    pregunta_usuario: str,
    clarificador_chain,
    sql_agent,
    explicador_chain,
    clasificador_chain,
    reformulador_chain,
    max_intentos=3,
):
    historial = []
    respuestas_usuario = {}
    prompt_actual = pregunta_usuario
    intentos = 0
    aclaraciones_str = ""

    while intentos < max_intentos:
        print(f"\n🔄 Iteración #{intentos + 1} - Refinando la pregunta...\n")

        # Paso 1: Clarificación guiada
        preguntas = clarificador_chain.run({"pregunta": prompt_actual}).strip()

        if preguntas == "NO_CLARIFICATION_NEEDED":
            print("✅ No hace falta pedir más aclaraciones.")
            prompt_claro = pregunta_usuario
        else:
            nuevas_respuestas = {}
            for pregunta in preguntas.split("\n"):
                if pregunta.strip():
                    user_input = input(f"{pregunta.strip()} 👉 ")
                    nuevas_respuestas[pregunta.strip()] = user_input

            respuestas_usuario.update(nuevas_respuestas)

            aclaraciones_str = "\n".join(
                f"- {k}: {v}" for k, v in respuestas_usuario.items()
            )
            prompt_claro = f"""Pregunta original: {pregunta_usuario}
        Aclaraciones:
        {aclaraciones_str}"""

        # Paso 2: Ejecutar la consulta SQL
        print("\n🤖 Ejecutando consulta...\n")
        try:
            resultado = sql_agent.run(prompt_claro)
        except Exception as e:
            error_msg = str(e)

            # Llamás a la chain de corrección
            fixed_query = correction_chain.run(
                {
                    "query": prompt_claro,
                    "error": error_msg,
                }
            )
            resultado = sql_agent.run(fixed_query.strip())

        if "error" in resultado.lower() or resultado.strip() == "":
            print("⚠️ La consulta no fue exitosa. Vamos a pedir más detalles...")
            prompt_actual = prompt_claro
            intentos += 1
            continue

        # Paso 3: Explicar el resultado
        explicacion = explicador_chain.run(
            {
                "pregunta": pregunta_usuario,
                "aclaraciones": aclaraciones_str,
                "resultado": resultado,
            }
        )

        print("\n🧠 Explicación final para el usuario:\n")
        print(explicacion)

        # Paso 4: Feedback
        feedback = input(
            "\n"
            + explicacion
            + "\n✍️ ¿Te resultó útil esta explicación? Podés responder con una frase 👉 "
        ).strip()
        clasificacion = (
            clasificador_chain.run({"respuesta_usuario": feedback}).strip().lower()
        )

        historial.append(
            {
                "pregunta": pregunta_usuario,
                "aclaraciones": respuestas_usuario.copy(),
                "prompt_final": prompt_claro,
                "resultado_sql": resultado,
                "explicacion": explicacion,
                "feedback_usuario": feedback,
                "clasificacion_feedback": clasificacion,
            }
        )

        if "útil" == clasificacion:
            print("✅ ¡Gracias! Me alegra que te haya servido.")
            return resultado, historial

        # Paso 5: Reformulación si no fue útil
        print("🔁 Gracias por tu comentario. Vamos a intentar mejorar la consulta...")

        contexto_historial = ""
        for h in historial:
            contexto_historial += f"""
[Pregunta anterior]: {h['pregunta']}
[Aclaraciones]: {h['aclaraciones']}
[Respuesta SQL]: {h['resultado_sql']}
[Explicación]: {h['explicacion']}
[Feedback]: {h['feedback_usuario']}
"""

        nueva_pregunta = reformulador_chain.run(
            {
                "historial": contexto_historial,
                "nueva_aclaracion": feedback,
                "pregunta_original": pregunta_usuario,
            }
        ).strip()

        print(f"\n📌 Reformulando la pregunta como:\n{nueva_pregunta}\n")
        prompt_actual = nueva_pregunta
        intentos += 1

    print("\n❌ No pudimos entender bien tu consulta después de varios intentos.")
    return None, historial


if __name__ == "__main__":
    pregunta = input("🧑‍💻 ¿Qué consulta querés hacer? 👉 ")
    resultado, historial = loop_consulta_sql(
        pregunta,
        clarificador_chain,
        agente,
        explicador_chain,
        clasificador_chain,
        reformulador_chain,
    )
    print("\n📊 Resultado de la consulta SQL:")
    if resultado:
        print(resultado)
        print("explicacion:", historial[-1]["explicacion"])
    else:
        print("⚠️ No se pudo obtener un resultado válido.")

    print("\n📜 Historial de la sesión:")
    for paso in historial:
        print(paso)
