import { Link } from 'react-router-dom';
import './Navigation.css';

export default function Navigation() {
  return (
    <nav className="main-nav">
      <ul>
        <li><Link to="/word-editor"> 修改词库</Link></li>
        <li><Link to="/word-search"> 词库搜索</Link></li>
        <li><Link to="/"> 返回主页</Link></li>
        <li class="right-item"><Link to="/"> 用户名</Link></li>
      </ul>
    </nav>
  );
}