from app.internal.config import get_settings
from app.internal.cache_manager import AgentCardCacheManager
from fastapi import FastAPI
from app.routers.health import router as health_router
from app.routers.agent_selection import router as agent_selection_router
from contextlib import asynccontextmanager

@asynccontextmanager
async def lifespan(app: FastAPI):
    settings = get_settings()
    cacheMgr = AgentCardCacheManager(settings)
    await cacheMgr.replace_cache()
    app.state.cacheMgr:AgentCardCacheManager = cacheMgr
    yield



app = FastAPI(title="Router", version="0.1.0", lifespan=lifespan)

app.include_router(health_router)
app.include_router(agent_selection_router)


