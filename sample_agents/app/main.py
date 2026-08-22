from anthropic import BaseModel
from fastapi import FastAPI
from app.music_agent import MusicAgent, Response
    
app = FastAPI(
    title="Sample Agents",
    description="Sample agents for demonstration purposes",
    version="0.0.1",
)

agent = MusicAgent() 

class Payload(BaseModel):
    message: str


@app.post("/music_agent", response_model=Response)
async def music_agent_endpoint(payload: Payload):
    return await agent.execute(payload.message)
