import httpx
from typing import Any
from app.internal.models import AgentCard

class KomikClient:
    def __init__(self, hostname: str):
        self.hostname = hostname

    async def get_agent_cards(self) -> list[AgentCard]:
        async with httpx.AsyncClient(base_url=self.hostname) as client:
            response = await client.get("/internal/api/v1/agent")
            response.raise_for_status()
            response_json = response.json()
            agents = [AgentCard(**agent) for agent in response_json["agents"]]
            return agents