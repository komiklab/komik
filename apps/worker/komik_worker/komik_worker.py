from temporalio.worker import Worker
from loguru import logger
from models.temporal_config import Configuration
from temporalio.client import Client
from komik_workflows.transition_workflow import EntityTransitionWorkflow
from shared.interceptors import LoguruInterceptor

class KomikWorker:
    def __init__(self):
        self.logger = logger
        self.config = Configuration()

    async def run(self):
        self.logger.info("Starting worker")
        client = await Client.connect(
            target_host=self.config.url,
            namespace=self.config.namespace
        )
        worker = Worker(
            client,
            task_queue=self.config.task_queue,
            workflows=[
                EntityTransitionWorkflow
            ],
            activities=[
            ],
            interceptors=[LoguruInterceptor()],
        )
        await worker.run()