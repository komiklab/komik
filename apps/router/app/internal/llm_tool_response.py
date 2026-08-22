import needle
from typing import List,Dict
import json


class LLMResponse:
    def __init__(self, tools:List[Dict]):
#         tools =[
#     {
#         "name": "music_agent",
#         "description": "Controls music playback. Use when the user asks to play, stop, pause, resume, or change music.",
#         "parameters": {
#             "type": "object",
#             "properties": {},
#             "required": []
#         }
#     },
#     {
#         "name": "thermostat_agent",
#         "description": "Controls the thermostat temperature. Use only when the user explicitly asks to set, increase, decrease, raise, lower, or change the temperature.",
#         "parameters": {
#             "type": "object",
#             "properties": {
#                 "temp": {
#                     "type": "float",
#                     "description": "Temperature to set"
#                 }
#             },
#             "required": ["temp"]
#         }
#     }
# ]
#         print(tools)
        self.instance = needle.Needle(tools=tools)

    def plan_tool_calls(self, message: str, max_steps: int = 8, threshold: float = 0.5) -> List[Dict]:

        self.instance.reset()
        plan = []
        response = self.instance.complete(text=message, max_new_tokens=256)

        for step in range(max_steps):
            calls = response.get("function_calls", [])
            confidence = response.get("confidence")
            ungrounded = response.get("validation", {}).get("ungrounded", [])
            if response.get("type") != "call" or not calls:
                return {
                    "function_calls": [p["call"] for p in plan],
                    "steps": plan,          # each accepted call WITH its own confidence
                    "status": "completed",
                    "steps_taken": step,
            }
            # if response.get("type") != "call" or not calls:
            #     return {
            #     "function_calls": plan,
            #     "status": "completed",
            #     "steps_taken": step,
            #     "confidence": confidence,
            #     "ungrounded": ungrounded,
            #     }  # "respond" with [] -> loop is done
            if confidence is None:
                status = "uncalibrated"
            elif confidence < threshold:
                status = "low_confidence"
            elif ungrounded:
                status = "ungrounded_field"
            else:
                status = "accepted"
            if status != "accepted":
                return {
                "function_calls": [p["call"] for p in plan],
                "steps": plan,
                "pending_call": calls,
                "confidence": confidence,
                "ungrounded": ungrounded,
                "status": status,
                "steps_taken": step,
            }

            for c in calls:
                plan.append({"call": c, "confidence": confidence, "ungrounded": ungrounded})

            # if status != "accepted":
            #     return {
            #     "function_calls": plan,          # steps accepted so far
            #     "pending_call": calls,           # the call that failed the gate
            #     "confidence": confidence,
            #     "ungrounded": ungrounded,
            #     "status": status,
            #     "steps_taken": step,
            #     }
            # plan.extend(calls)
            # name = calls[0]["name"]
            response = self.instance.complete(text=json.dumps({}), max_new_tokens=256)

        return {
        "function_calls": plan,
        "status": "max_steps_reached",
        "steps_taken": max_steps,
    }
            
