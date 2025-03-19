import { Link } from 'react-router-dom';

export default function Navigation() {
  return (
    <nav className="main-nav">
      <ul>
        <li><Link to="/word-editor">📖 修改词库</Link></li>
        <li><Link to="/">🏠 返回主页</Link></li>
      </ul>
    </nav>
  );
}