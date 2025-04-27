import React from 'react';
import './AboutUs.css';

const AboutUs = () => {
  return (
    <div className="about-page">
      <header className="main-header">
        <h1>关于我们</h1>
        <h2>About Us</h2>
      </header>

      <main className="about-content">
        <section className="about-section">
          <h3>项目说明</h3>
          <p>本项目的目的在于为日语学习者提供一个简洁、易于使用的网页，主要注重于语法、单词的复习。于是
            根据作者本人的日语学习经历，通过赋予权重的形式，以更符合记忆曲线的方式提供复习。这是我认为现有的单词学习类应用所不具备的。
          </p>
          <p>
            当然，本项目并非成熟的商业级项目，只是作者在学习日语之余开发的小玩具，如果能帮上您，那再好不过。
          </p>
        </section>

        <section className="about-section">
          <h3>如何使用</h3>
          <p>用户可以自己部署该项目(本地即可，也可以部署在云，和其他朋友一起使用)。目前数据库还未准备完善，故用户只能自己边学习，边记录自己所学。而后端回为用户提供复习服务。
            用户同时也可以在自己的数据库中进行检索和修改。
          </p>
          <p>

          </p>
        </section>


        <section className="about-section contact-section">
          <h3>联系我们</h3>
            <p>如有任何问题或建议，欢迎通过以下方式联系我们：</p>
            <ul>
              <li>邮箱: <a href="mailto:sunqizhen6@gmail.com">sunqizhen6@gmail.com</a></li>
              <li>GitHub: <a href="https://github.com/lamperriat/easy-japanese" target="_blank" rel="noopener noreferrer">easy-japanese</a></li>
            </ul>
        </section>
      </main>
    </div>
  );
};

export default AboutUs;