import asyncio
import operator
from random import randint
from typing import Annotated, AsyncGenerator, TypedDict, Any
from langgraph.graph import StateGraph, START, END

class CVEState(TypedDict, total=False):
    cve_id: str
    status: str
    worktree_path: str
    updated_dependencies: list[str]
    test_results: str
    pr_raised: str
    logs: Annotated[list[str], operator.add]

async def create_git_worktree(state: CVEState) -> dict:
    cve_id = state.get("cve_id", "CVE-UNKNOWN")
    worktree = f"worktree/fix-{cve_id.lower().replace('-', '_')}"
    return {
        "status": "worktree_created",
        "worktree_path": worktree,
        "logs": [f"Git worktree created at {worktree}"],
    }

async def update_go_mod(state: CVEState) -> dict:
    cve_id = state.get("cve_id", "CVE-UNKNOWN")
    deps = ["golang.org/x/crypto@v0.17.0", "github.com/gin-gonic/gin@v1.9.1"]
    return {
        "status": "go_mod_updated",
        "updated_dependencies": deps,
        "logs": [f"go.mod file updated with patched dependencies for {cve_id}"],
    }

async def run_unit_tests(state: CVEState) -> dict:
    return {
        "status": "tests_passed",
        "test_results": "PASS (42 tests run, 0 failures)",
        "logs": ["Unit tests passed successfully (42/42)"],
    }

async def raise_pull_request(state: CVEState) -> dict:
    cve_id = state.get("cve_id", "CVE-UNKNOWN")
    pr_num = randint(100, 999)
    pr_ref = f"prref{pr_num}"
    return {
        "status": "completed",
        "pr_raised": pr_ref,
        "logs": [f"Pull request created: {pr_ref}"],
    }

builder = StateGraph(CVEState)
builder.add_node("create_git_worktree", create_git_worktree)
builder.add_node("update_go_mod", update_go_mod)
builder.add_node("run_unit_tests", run_unit_tests)
builder.add_node("raise_pull_request", raise_pull_request)

builder.add_edge(START, "create_git_worktree")
builder.add_edge("create_git_worktree", "update_go_mod")
builder.add_edge("update_go_mod", "run_unit_tests")
builder.add_edge("run_unit_tests", "raise_pull_request")
builder.add_edge("raise_pull_request", END)

graph = builder.compile()

async def stream_cve_workflow(cve_id: str) -> AsyncGenerator[dict[str, Any], None]:
    async for event in graph.astream({"cve_id": cve_id}, stream_mode="updates"):
        for node_name, output in event.items():
            yield {
                "node": node_name,
                "status": output.get("status"),
                "logs": output.get("logs", []),
                "data": output,
            }
