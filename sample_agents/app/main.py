import json
from fastapi import FastAPI
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

from app.music_agent import MusicAgent, Response
from app.cve_workflow import graph, stream_cve_workflow
    
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

class CVERequest(BaseModel):
    cve_id: str

@app.post("/workflows/cve")
async def run_cve_workflow(request: CVERequest):
    result = await graph.ainvoke({"cve_id": request.cve_id})
    return result

@app.post("/workflows/cve/stream")
async def stream_cve_endpoint(request: CVERequest):
    async def event_generator():
        async for event in stream_cve_workflow(request.cve_id):
            yield f"data: {json.dumps(event)}\n\n"

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
        },
    )


