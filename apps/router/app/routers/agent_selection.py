from app.internal.models import convertAgentCardtoNeedleAcceptedCar
from app.internal.llm_tool_response import LLMResponse
from app.internal.models import AgentCard
from app.internal.config import get_settings
from fastapi import Depends
from app.internal.config import Settings
from fastapi import APIRouter, Request
from typing import Annotated

router = APIRouter(prefix="/agent-selection", tags=["agent-selection"])


@router.get("/")
async def agent_selection(request: Request, message: str, settings:Annotated[Settings, Depends(get_settings)]):
    cache_manager:list[AgentCard] = request.app.state.cacheMgr.agentCardCache
    # Convert AgentCard to NeedleAcceptedCard
    needle_accepted_cards = []
    for agent in cache_manager:
        needle_accepted_cards.append(convertAgentCardtoNeedleAcceptedCar(agent))
    _llm = LLMResponse(tools=[agent.dict() for agent in needle_accepted_cards])
    response = _llm.plan_tool_calls(message)
    return response