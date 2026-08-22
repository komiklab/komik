import asyncio
from enum import Enum
from pydantic_ai import Agent
from pydantic_ai.models.test import TestModel
from pydantic import BaseModel, Field

class StatusEnum(str, Enum):
    on = "on"
    off = "off"

class MusicTypeEnum(str, Enum):
    classical = "classical"
    jazz = "jazz"
    rock = "rock"
    pop = "pop"
    
class Response(BaseModel):
    # status must be enum either on or off
    status:StatusEnum
    # if status is on then type_music must be enum either classical or jazz or rock or pop
    type_music:MusicTypeEnum
    volume:int = Field(description="volume must be between 0 and 100")

class Input(BaseModel):
    command:str


class MusicAgent:
    def __init__(self):
        self.agent = Agent(
            model = TestModel(),
            output_type=Response,
            system_prompt="""
            You are a music agent. You will be given a command and you will return the status of the music player. 
            The response should be in JSON format and should have the following fields:
            - status: enum of on or off
            - type_music: enum of classical or jazz or rock or pop
            - volume: int between 0 and 100
            """
        )

    async def execute(self, command:str):
        response = await self.agent.run(command)
        return response.output

if __name__ == "__main__":
    music_agent = MusicAgent()
    asyncio.run(music_agent.execute("turn on the music"))
