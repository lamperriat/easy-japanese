import React from 'react';
import './HomePage.css' 
import pic1 from './assets/homepage/pic1.jpg';
import pic2 from './assets/homepage/pic2.jpg';
import pic3 from './assets/homepage/pic3.jpg';
export default function HomePage() {
  return (
    <div className="home-page">
      <header className="main-header">
        <h1>易学日语</h1>
        <h2>やさしい日本語</h2>
      </header>
      
      <main className="home-content">
        <section className="welcome-section">
          <h3>欢迎使用！</h3>
          <p>选择导航栏开始学习！</p>
        </section>

        <div className="image-section">
          <img className="home-image-1" src={pic1} alt="image1" />
          <img className="home-image-2" src={pic2} alt="image2" />
          <img className="home-image-3" src={pic3} alt="image3" />
        </div>
      </main>
    </div>
  );
}
