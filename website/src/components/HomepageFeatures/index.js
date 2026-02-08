import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'OpenAI Compatible',
    emoji: '🚀',
    description: (
      <>
        Seamless integration with OpenAI SDK and compatible clients. 
        Full support for chat completions, embeddings, and streaming responses.
      </>
    ),
  },
  {
    title: 'Intelligent Routing',
    emoji: '🎯',
    description: (
      <>
        Multiple routing strategies including round-robin, random, weighted, 
        model-based matching, and failover for optimal request distribution.
      </>
    ),
  },
  {
    title: 'Easy Management',
    emoji: '📊',
    description: (
      <>
        Modern web-based admin dashboard with comprehensive management features 
        for LLM resources, API keys, policies, and system monitoring.
      </>
    ),
  },
  {
    title: 'Secure & Flexible',
    emoji: '🔐',
    description: (
      <>
        Flexible authentication system with API key management, API key support, 
        and configurable security settings.
      </>
    ),
  },
  {
    title: 'Production Ready',
    emoji: '⚡',
    description: (
      <>
        Built with Go and Vue 3, supports multiple storage backends, 
        Docker deployment, and comprehensive logging.
      </>
    ),
  },
  {
    title: 'Multi-Language',
    emoji: '🌐',
    description: (
      <>
        Full internationalization support with Chinese and English interfaces. 
        Complete documentation in multiple languages.
      </>
    ),
  },
];

function Feature({emoji, title, description}) {
  return (
    <div className={clsx('col col--4', styles.featureCard)}>
      <div className="text--center">
        <div className={styles.featureEmoji}>{emoji}</div>
      </div>
      <div className="text--center padding-horiz--md">
        <h3 className={styles.featureTitle}>{title}</h3>
        <p className={styles.featureDescription}>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className={styles.sectionHeader}>
          <h2 className={styles.sectionTitle}>Core Features</h2>
          <p className={styles.sectionSubtitle}>
            Complete intelligent API gateway solution for AI applications
          </p>
        </div>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
