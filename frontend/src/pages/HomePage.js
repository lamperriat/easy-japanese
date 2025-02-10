import React from 'react';

export default function HomePage() {
  return (
    <div className="home-page">
      <header className="main-header">
        <h1>易学日语</h1>
        <h2>Easy Japanese Learning</h2>
      </header>
      
      <main className="home-content">
        <section className="welcome-section">
          <h3>欢迎使用！</h3>
          <p>选择上方导航栏开始学习之旅</p>
        </section>
      </main>
    </div>
  );
}