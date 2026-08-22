from fastapi import FastAPI
from app.music_agent import MusicAgent, Response

app = FastAPI(
    title="Sample Agents",
    description="Sample agents for demonstration purposes",
    version="0.0.1",
)

agent = MusicAgent() 

@app.get("/music_agent", response_model=Response)
async def music_agent_endpoint(input: str):
    return await agent.execute(input)

