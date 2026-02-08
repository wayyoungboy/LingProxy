import React, { useState } from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const codeExamples = [
  {
    language: 'python',
    label: 'Python',
    code: `from openai import OpenAI

client = OpenAI(
    api_key="your-token",
    base_url="http://localhost:8080/llm/v1"
)

response = client.chat.completions.create(
    model="gpt-3.5-turbo",
    messages=[
        {"role": "user", "content": "Hello!"}
    ]
)

print(response.choices[0].message.content)`,
  },
  {
    language: 'javascript',
    label: 'JavaScript',
    code: `import OpenAI from 'openai';

const client = new OpenAI({
  apiKey: 'your-token',
  baseURL: 'http://localhost:8080/llm/v1'
});

const completion = await client.chat.completions.create({
  model: 'gpt-3.5-turbo',
  messages: [{ role: 'user', content: 'Hello!' }]
});

console.log(completion.choices[0].message.content);`,
  },
  {
    language: 'bash',
    label: 'cURL',
    code: `curl -X POST http://localhost:8080/llm/v1/chat/completions \\
  -H "Authorization: Bearer your-token" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ]
  }'`,
  },
];

export default function CodeExample() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <div className={styles.codeExample}>
      <div className={styles.codeTabs}>
        {codeExamples.map((example, index) => (
          <button
            key={index}
            className={clsx(styles.tab, {
              [styles.activeTab]: activeTab === index,
            })}
            onClick={() => setActiveTab(index)}>
            {example.label}
          </button>
        ))}
      </div>
      <div className={styles.codeBlock}>
        <pre className={styles.code}>
          <code className={`language-${codeExamples[activeTab].language}`}>
            {codeExamples[activeTab].code}
          </code>
        </pre>
      </div>
    </div>
  );
}
