import React from 'react';
import './HomePage.css'; // 引入样式文件

export default function HomePage() {
  return (
    <div className="home-page">
      <header className="main-header">
        <h1>易学日语</h1>
        <h2>やさしい日本語</h2>
      </header>
      
      <main className="home-content">
        {/* 文本居中部分 */}
        <section className="welcome-section">
          <h3>欢迎使用！</h3>
          <p>选择导航栏开始学习！</p>
        </section>

        {/* 图片部分，左边和右边各一张 */}
        <div className="image-section">
          <img className="home-image-1" src="/pic/pic1.jpg" alt="image1" />
          <img className="home-image-2" src="/pic/pic2.jpg" alt="image2" />
          <img className="home-image-3" src="/pic/pic3.jpg" alt="image3" />
        </div>
      </main>
    </div>
  );
}
