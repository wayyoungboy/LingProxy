import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import CodeExample from '@site/src/components/CodeExample';

import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className={styles.heroBackground}></div>
      <div className="container">
        <div className={styles.heroContent}>
          <h1 className={clsx('hero__title', styles.heroTitle)}>
            {siteConfig.title}
          </h1>
          <p className={clsx('hero__subtitle', styles.heroSubtitle)}>
            {siteConfig.tagline}
          </p>
          <p className={styles.heroDescription}>
            Get started in minutes, scale to millions
          </p>
          <div className={styles.buttons}>
            <Link
              className={clsx('button button--primary button--lg', styles.ctaButton)}
              to="/docs/introduction">
              Get Started →
            </Link>
            <Link
              className={clsx('button button--outline button--secondary button--lg', styles.secondaryButton)}
              to="/docs/quick-start">
              View Docs
            </Link>
          </div>
        </div>
        <CodeExample />
      </div>
    </header>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`${siteConfig.title} - ${siteConfig.tagline}`}
      description="High-performance AI API Gateway with OpenAI-compatible interfaces, intelligent load balancing, and comprehensive management features.">
      <HomepageHeader />
      <main>
        <HomepageFeatures />
      </main>
    </Layout>
  );
}
