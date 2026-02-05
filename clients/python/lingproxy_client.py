#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
LingProxy Python Client
A standard Python client for LingProxy AI API Gateway

Usage:
    from lingproxy_client import LingProxyClient
    
    client = LingProxyClient(
        api_key="your-api-key",
        base_url="http://localhost:8080/llm/v1"
    )
    
    response = client.chat.completions.create(
        model="gpt-3.5-turbo",
        messages=[
            {"role": "user", "content": "Hello!"}
        ]
    )
"""

import os
from typing import Optional, List, Dict, Any, Union
from openai import OpenAI


class LingProxyClient:
    """
    LingProxy Python Client
    
    A standard Python client for interacting with LingProxy AI API Gateway.
    Compatible with OpenAI Python SDK.
    """
    
    def __init__(
        self,
        api_key: Optional[str] = None,
        base_url: str = "http://localhost:8080/llm/v1",
        timeout: float = 30.0,
    ):
        """
        Initialize LingProxy client
        
        Args:
            api_key: API key or token. If not provided, will try to get from 
                    LINGPROXY_API_KEY environment variable
            base_url: Base URL of LingProxy server (default: http://localhost:8080/llm/v1)
            timeout: Request timeout in seconds (default: 30.0)
        """
        if api_key is None:
            api_key = os.getenv("LINGPROXY_API_KEY")
            if api_key is None:
                raise ValueError(
                    "API key is required. Either pass it as parameter or set "
                    "LINGPROXY_API_KEY environment variable."
                )
        
        self._client = OpenAI(
            api_key=api_key,
            base_url=base_url,
            timeout=timeout,
        )
        
        # Expose OpenAI client's attributes for compatibility
        self.chat = self._client.chat
        self.completions = self._client.completions
        self.models = self._client.models
        self.embeddings = self._client.embeddings
    
    def list_models(self) -> List[Dict[str, Any]]:
        """
        List all available models
        
        Returns:
            List of model dictionaries
        """
        models = self._client.models.list()
        return [model.model_dump() for model in models.data]
    
    def create_chat_completion(
        self,
        model: str,
        messages: List[Dict[str, str]],
        temperature: Optional[float] = None,
        max_tokens: Optional[int] = None,
        top_p: Optional[float] = None,
        stream: bool = False,
        **kwargs
    ) -> Dict[str, Any]:
        """
        Create a chat completion
        
        Args:
            model: Model name to use
            messages: List of message dictionaries with 'role' and 'content'
            temperature: Sampling temperature (0-2)
            max_tokens: Maximum tokens to generate
            top_p: Nucleus sampling parameter
            stream: Whether to stream the response
            **kwargs: Additional parameters
            
        Returns:
            Chat completion response dictionary
        """
        params = {
            "model": model,
            "messages": messages,
            "stream": stream,
        }
        
        if temperature is not None:
            params["temperature"] = temperature
        if max_tokens is not None:
            params["max_tokens"] = max_tokens
        if top_p is not None:
            params["top_p"] = top_p
        
        params.update(kwargs)
        
        response = self._client.chat.completions.create(**params)
        
        if stream:
            return response
        
        return {
            "id": response.id,
            "object": response.object,
            "created": response.created,
            "model": response.model,
            "choices": [
                {
                    "index": choice.index,
                    "message": {
                        "role": choice.message.role,
                        "content": choice.message.content,
                    },
                    "finish_reason": choice.finish_reason,
                }
                for choice in response.choices
            ],
            "usage": {
                "prompt_tokens": response.usage.prompt_tokens,
                "completion_tokens": response.usage.completion_tokens,
                "total_tokens": response.usage.total_tokens,
            },
        }
    
    def create_completion(
        self,
        model: str,
        prompt: str,
        temperature: Optional[float] = None,
        max_tokens: Optional[int] = None,
        **kwargs
    ) -> Dict[str, Any]:
        """
        Create a text completion
        
        Args:
            model: Model name to use
            prompt: Text prompt
            temperature: Sampling temperature (0-2)
            max_tokens: Maximum tokens to generate
            **kwargs: Additional parameters
            
        Returns:
            Completion response dictionary
        """
        params = {
            "model": model,
            "prompt": prompt,
        }
        
        if temperature is not None:
            params["temperature"] = temperature
        if max_tokens is not None:
            params["max_tokens"] = max_tokens
        
        params.update(kwargs)
        
        response = self._client.completions.create(**params)
        
        return {
            "id": response.id,
            "object": response.object,
            "created": response.created,
            "model": response.model,
            "choices": [
                {
                    "index": choice.index,
                    "text": choice.text,
                    "finish_reason": choice.finish_reason,
                }
                for choice in response.choices
            ],
            "usage": {
                "prompt_tokens": response.usage.prompt_tokens,
                "completion_tokens": response.usage.completion_tokens,
                "total_tokens": response.usage.total_tokens,
            },
        }


# Example usage
if __name__ == "__main__":
    import json
    
    # Initialize client
    # Option 1: Use environment variable
    # export LINGPROXY_API_KEY=your-api-key
    # client = LingProxyClient()
    
    # Option 2: Pass API key directly
    api_key = os.getenv("LINGPROXY_API_KEY", "ling-Uc9tFvNr97HaMXrKXw2R1ZqNiHt_pp0M_OsDOvjns8M=")
    client = LingProxyClient(
        api_key=api_key,
        base_url="http://localhost:8080/llm/v1"
    )
    
    print("=" * 70)
    print("LingProxy Python Client Demo")
    print("=" * 70)
    print()
    
    # Example 1: List models
    print("Example 1: List available models")
    print("-" * 70)
    try:
        models = client.list_models()
        print(f"Found {len(models)} models:")
        for model in models[:5]:  # Show first 5
            print(f"  - {model.get('id', 'N/A')}")
    except Exception as e:
        print(f"Error: {e}")
    print()
    
    # Example 2: Chat completion
    print("Example 2: Chat completion")
    print("-" * 70)
    try:
        response = client.create_chat_completion(
            model="glm-4.5-flash",  # Replace with your model name
            messages=[
                {"role": "system", "content": "You are a helpful assistant."},
                {"role": "user", "content": "Tell me about Greece's largest city."}
            ],
            temperature=0.7,
            max_tokens=100
        )
        print(f"Response ID: {response['id']}")
        print(f"Model: {response['model']}")
        if response['choices']:
            content = response['choices'][0]['message']['content']
            print(f"Content: {content[:200]}..." if len(content) > 200 else f"Content: {content}")
        print(f"Usage: {response['usage']}")
    except Exception as e:
        print(f"Error: {e}")
    print()
    
    # Example 3: Using OpenAI SDK style (direct access)
    print("Example 3: Using OpenAI SDK style")
    print("-" * 70)
    try:
        response = client.chat.completions.create(
            model="glm-4.5-flash",
            messages=[
                {"role": "user", "content": "Hello! How are you?"}
            ]
        )
        print(f"Response: {response.choices[0].message.content}")
    except Exception as e:
        print(f"Error: {e}")
    print()
    
    print("=" * 70)
    print("Demo completed!")
    print("=" * 70)
