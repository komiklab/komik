from app.internal.models import AgentCard
from app.internal.config import Settings
from app.internal.komik_client import KomikClient

class AgentCardCacheManager:
    def __init__(self, settings: Settings)->None:
        self.cache: list[AgentCard] = []
        self.komik_client = KomikClient(settings.komik_endpoint)

    async def _populate_cache(self)->None:
        try:
            self.cache = await self.komik_client.get_agent_cards()
        except Exception as e:
            raise

    async def replace_cache(self)->None:
        await self._populate_cache()

    @property
    def agentCardCache(self)->list[AgentCard]:
        return self.cache
    
    

