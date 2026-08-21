from pydantic import Field
from pydantic import BaseModel
from typing import Optional

class Parameter(BaseModel):
    name: str
    description: str
    type: str
    required: bool

class AgentCard(BaseModel):
    id: str
    description: str
    name: str
    capabilities: list[str]
    parameter: Optional[list[dict]] = []

class NeedleAcceptedCard(BaseModel):
    name: str
    description: str
    parameters: NeedleParameters

class NeedleParameters(BaseModel):
    type: str = "object"
    properties: dict[str, NeedleParameterProperties] = Field(default_factory=dict)
    required: list[str]= Field(default_factory=list)

class NeedleParameterProperties(BaseModel):
    type: str
    description: str
    

def convertAgentCardtoNeedleAcceptedCar(ac:AgentCard) -> NeedleAcceptedCard:
    properties = {
        param["name"]: NeedleParameterProperties(
            type=param["type"],
            description=param["description"],
        )
        for param in ac.parameter
    }
    required = [
        param["name"]
        for param in ac.parameter
        if param["required"]
    ]

    return NeedleAcceptedCard(
        name=ac.name,
        description=ac.description,
        parameters=NeedleParameters(
            type="object",
            properties=properties,
            required=required,
        ),
    )

    