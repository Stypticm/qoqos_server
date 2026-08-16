from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from agent import process_query
import traceback

app = FastAPI(title="AI-Agent для ремонта телефонов")

class QueryRequest(BaseModel):
    user_query: str

@app.post("/ask")
async def ask_agent(request: QueryRequest):
    try:
        answer = await process_query(request.user_query)
        return {"answer": answer}
    except Exception as e:
        print(traceback.format_exc())
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/ping")
async def ping():
    return {"status": "ok", "message": "AI-Agent is running with qwen:4b"}