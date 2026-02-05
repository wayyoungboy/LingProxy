/**
 * LingProxy JavaScript Client
 * A standard JavaScript/TypeScript client for LingProxy AI API Gateway
 * 
 * Usage:
 *   import { LingProxyClient } from './lingproxy-client.js';
 *   
 *   const client = new LingProxyClient({
 *     apiKey: 'your-api-key',
 *     baseURL: 'http://localhost:8080/llm/v1'
 *   });
 *   
 *   const response = await client.chat.completions.create({
 *     model: 'gpt-3.5-turbo',
 *     messages: [
 *       { role: 'user', content: 'Hello!' }
 *     ]
 *   });
 */

import OpenAI from 'openai';

export class LingProxyClient {
  /**
   * Initialize LingProxy client
   * 
   * @param {Object} options - Client options
   * @param {string} options.apiKey - API key or token. If not provided, will try to get from LINGPROXY_API_KEY environment variable
   * @param {string} options.baseURL - Base URL of LingProxy server (default: http://localhost:8080/llm/v1)
   * @param {number} options.timeout - Request timeout in milliseconds (default: 30000)
   */
  constructor(options = {}) {
    const {
      apiKey = process.env.LINGPROXY_API_KEY || process.env.OPENAI_API_KEY,
      baseURL = 'http://localhost:8080/llm/v1',
      timeout = 30000,
    } = options;

    if (!apiKey) {
      throw new Error(
        'API key is required. Either pass it as parameter or set ' +
        'LINGPROXY_API_KEY environment variable.'
      );
    }

    this._client = new OpenAI({
      apiKey,
      baseURL,
      timeout,
    });

    // Expose OpenAI client's attributes for compatibility
    this.chat = this._client.chat;
    this.completions = this._client.completions;
    this.models = this._client.models;
    this.embeddings = this._client.embeddings;
  }

  /**
   * List all available models
   * 
   * @returns {Promise<Array>} List of model objects
   */
  async listModels() {
    const models = await this._client.models.list();
    return models.data.map(model => ({
      id: model.id,
      object: model.object,
      created: model.created,
      owned_by: model.owned_by,
    }));
  }

  /**
   * Create a chat completion
   * 
   * @param {Object} params - Chat completion parameters
   * @param {string} params.model - Model name to use
   * @param {Array} params.messages - List of message objects with 'role' and 'content'
   * @param {number} [params.temperature] - Sampling temperature (0-2)
   * @param {number} [params.max_tokens] - Maximum tokens to generate
   * @param {number} [params.top_p] - Nucleus sampling parameter
   * @param {boolean} [params.stream] - Whether to stream the response
   * @returns {Promise<Object>} Chat completion response
   */
  async createChatCompletion(params) {
    const {
      model,
      messages,
      temperature,
      max_tokens,
      top_p,
      stream = false,
      ...rest
    } = params;

    const requestParams = {
      model,
      messages,
      stream,
      ...rest,
    };

    if (temperature !== undefined) requestParams.temperature = temperature;
    if (max_tokens !== undefined) requestParams.max_tokens = max_tokens;
    if (top_p !== undefined) requestParams.top_p = top_p;

    const response = await this._client.chat.completions.create(requestParams);

    if (stream) {
      return response;
    }

    return {
      id: response.id,
      object: response.object,
      created: response.created,
      model: response.model,
      choices: response.choices.map(choice => ({
        index: choice.index,
        message: {
          role: choice.message.role,
          content: choice.message.content,
        },
        finish_reason: choice.finish_reason,
      })),
      usage: {
        prompt_tokens: response.usage.prompt_tokens,
        completion_tokens: response.usage.completion_tokens,
        total_tokens: response.usage.total_tokens,
      },
    };
  }

  /**
   * Create a text completion
   * 
   * @param {Object} params - Completion parameters
   * @param {string} params.model - Model name to use
   * @param {string} params.prompt - Text prompt
   * @param {number} [params.temperature] - Sampling temperature (0-2)
   * @param {number} [params.max_tokens] - Maximum tokens to generate
   * @returns {Promise<Object>} Completion response
   */
  async createCompletion(params) {
    const {
      model,
      prompt,
      temperature,
      max_tokens,
      ...rest
    } = params;

    const requestParams = {
      model,
      prompt,
      ...rest,
    };

    if (temperature !== undefined) requestParams.temperature = temperature;
    if (max_tokens !== undefined) requestParams.max_tokens = max_tokens;

    const response = await this._client.completions.create(requestParams);

    return {
      id: response.id,
      object: response.object,
      created: response.created,
      model: response.model,
      choices: response.choices.map(choice => ({
        index: choice.index,
        text: choice.text,
        finish_reason: choice.finish_reason,
      })),
      usage: {
        prompt_tokens: response.usage.prompt_tokens,
        completion_tokens: response.usage.completion_tokens,
        total_tokens: response.usage.total_tokens,
      },
    };
  }
}

// Example usage
if (import.meta.url === `file://${process.argv[1]}`) {
  (async () => {
    const apiKey = process.env.LINGPROXY_API_KEY || 'ling-Uc9tFvNr97HaMXrKXw2R1ZqNiHt_pp0M_OsDOvjns8M=';
    
    const client = new LingProxyClient({
      apiKey,
      baseURL: 'http://localhost:8080/llm/v1'
    });

    console.log('='.repeat(70));
    console.log('LingProxy JavaScript Client Demo');
    console.log('='.repeat(70));
    console.log();

    // Example 1: List models
    console.log('Example 1: List available models');
    console.log('-'.repeat(70));
    try {
      const models = await client.listModels();
      console.log(`Found ${models.length} models:`);
      models.slice(0, 5).forEach(model => {
        console.log(`  - ${model.id}`);
      });
    } catch (error) {
      console.error(`Error: ${error.message}`);
    }
    console.log();

    // Example 2: Chat completion
    console.log('Example 2: Chat completion');
    console.log('-'.repeat(70));
    try {
      const response = await client.createChatCompletion({
        model: 'glm-4.5-flash', // Replace with your model name
        messages: [
          { role: 'system', content: 'You are a helpful assistant.' },
          { role: 'user', content: 'Tell me about Greece\'s largest city.' }
        ],
        temperature: 0.7,
        max_tokens: 100
      });
      console.log(`Response ID: ${response.id}`);
      console.log(`Model: ${response.model}`);
      if (response.choices && response.choices.length > 0) {
        const content = response.choices[0].message.content;
        console.log(`Content: ${content.length > 200 ? content.substring(0, 200) + '...' : content}`);
      }
      console.log(`Usage:`, response.usage);
    } catch (error) {
      console.error(`Error: ${error.message}`);
    }
    console.log();

    // Example 3: Using OpenAI SDK style
    console.log('Example 3: Using OpenAI SDK style');
    console.log('-'.repeat(70));
    try {
      const response = await client.chat.completions.create({
        model: 'glm-4.5-flash',
        messages: [
          { role: 'user', content: 'Hello! How are you?' }
        ]
      });
      console.log(`Response: ${response.choices[0].message.content}`);
    } catch (error) {
      console.error(`Error: ${error.message}`);
    }
    console.log();

    console.log('='.repeat(70));
    console.log('Demo completed!');
    console.log('='.repeat(70));
  })();
}

export default LingProxyClient;
