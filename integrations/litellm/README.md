# Framework Integration

Resolve an approved prompt version before any model call, and keep its immutable prompt id and version beside the call for tracing. Put the resolved template in the system message using the framework's native message object.

LiteLLM metadata, LangChain runnable config, LangGraph state, and custom application tracing can all hold this information. This library is framework-agnostic and is not a model provider.