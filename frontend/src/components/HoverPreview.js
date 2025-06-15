
import React, { useState, useRef } from "react";
import "./HoverPreview.css";
const HoverPreview = (
  {text, fGetContent} 
) => {
  const [visible, setVisible] = useState(false);
  const [content, setContent] = useState("");

  const handleOnMouseEnter = () => {
    setVisible(true);
    if (!content) {
      fGetContent(text).then((data) => {
        if (data) {
          setContent(data);
        } else {
          setContent("<p>No content available </p>");
        }
      }).catch((error) => {
        setContent("<p>Error fetching content </p>");
      });
    }
    // console.log(content);
  }

  return (
    <div className="hover-container"
    // here we need to specify the file such that
    // the div can be inlined with the text
    style={{ display: 'inline-block' }}>
    <span 
      className="hover-trigger"
      onMouseEnter={handleOnMouseEnter}
      onMouseLeave={() => setVisible(false)}
    >{text}</span>
    
    {
      visible && content && (
        <div className="hover-content">
          <div dangerouslySetInnerHTML={{ __html: content }} />
        </div>
      )
    }
    </div>
  )
}

export default HoverPreview;