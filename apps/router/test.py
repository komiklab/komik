import needle

# @needle.tool
# def change_temperature(temperature: int):
#     """Change the temperature to a new value.

#     Args:
#         temperature: target temperature in degrees Celsius
#     """
#     return {"new_temperature": temperature}

tools = {
    "name": "change_temperature",
    "description": "Change the temperature of the thermostat to a specified target temperature.",
    "parameters": {
        "type": "object",
        "properties": {
            "temperature": {
                "type": "number",
                "description": "The target temperature in degrees Celsius."
            }
        },
        "required": ["temperature"]
    }
}

agent = needle.Needle(tools=[tools])
agent.reset()
result = agent.run("make the temperature 20 degrees")
print(result["results"])
