import React, { useState, useRef, useEffect, useCallback } from 'react';
import './VideoPlayer.css';
import subParse from '../components/SubParser/SubParser';
import { API_BASE_URL } from '../services/api';
import Notification from '../components/Auth/Notification';
import DictSearch from '../components/Dict/DictSearch';
import HoverPreview from '../components/HoverPreview';

const VideoMode = {
  WATCH: 'watch',
  STUDY: 'study',
  INVALID: 'invalid',
  [Symbol.for('isEnum')]: true
}

const VideoPlayer = () => {
  const [videoUrl, setVideoUrl] = useState('');
  const [isPlaying, setIsPlaying] = useState(false);
  const [volume, setVolume] = useState(1);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [showControls, setShowControls] = useState(true);
  const [fullscreen, setFullscreen] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');
  const [selectedFile, setSelectedFile] = useState(null);
  const [videoMode, setVideoMode] = useState(VideoMode.WATCH);
  const [subtitleUrl, setSubtitleUrl] = useState('');
  const rSubtitleLn = useRef(0);
  const [subtitleFileName, setSubtitleFileName] = useState('');
  const [curSubtitleLineIndex, setCurSubtitleLineIndex] = useState(0);
  const [assLines, setAssLines] = useState([]);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });

  const [prevJPLine, setPrevJPLine] = useState('');
  const [curJPLine, setCurJPLine] = useState('');
  const [nextJPLine, setNextJPLine] = useState('');

  // settings
  const [subtitleTimeOffset, setSubtitleTimeOffset] = useState(0);
  const [dictUrl, setDictUrl] = useState('https://dict.youdao.com/result?word={}&lang=ja');

  async function GetTokenizedString(sentence) {
    var token = sessionStorage.getItem('token');
    if (!token) {
      setNotification({
        show: true,
        message: '请先登录',
        type: 'error'
      });
      setTimeout(() => {
        setNotification({ show: false, message: '', type: '' });
      }, 3000);
      return;
    }
    try {
      const response = await fetch(`${API_BASE_URL}/api/subtitles/tokenize`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token,
        },
        body: JSON.stringify({ text: sentence })
      });
      const data = await response.json();
      if (response.ok && data && data.tokens) {
        return data.tokens;
      }
      console.log(data.error || 'Unknown error');
      return [];
    } catch (error) {
      console.error("Error fetching tokenized string:", error);
      setNotification({
        show: true,
        message: '网络请求失败',
        type: 'error'
      });
      setTimeout(() => {
        setNotification({ show: false, message: '', type: '' });
      }, 3000);
      return '';
    }
    
  }

  async function GetHTMLStringProxy(word) {
    var token = sessionStorage.getItem('token');
    if (!token) {
      setNotification({
        show: true,
        message: '请先登录',
        type: 'error'
      });
      setTimeout(() => {
        setNotification({ show: false, message: '', type: '' });
      }, 3000);
      return;
    }
    try {
      return await DictSearch(word, token);
    } catch (error) {
      return '无法找到内容，可能是网络原因，或查找的不是常规单词';
    }
  }

  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      // Reset error message and player state
      setErrorMessage('');
      setCurrentTime(0);
      setDuration(0);
      setIsPlaying(false);
      
      const url = URL.createObjectURL(file);
      setVideoUrl(url);
      setSelectedFile(file);
      
      // Show success notification
      setErrorMessage(`已选择文件: ${file.name}`);
      setTimeout(() => setErrorMessage(''), 3000);
      
      // Attempt to autoplay when ready
      setTimeout(() => {
        if (videoRef.current) {
          videoRef.current.play()
            .then(() => setIsPlaying(true))
            .catch(err => console.log("Autoplay prevented: ", err));
        }
      }, 500);
    }
  };

  const handleSubtitleFileSelect = (e) => {
    const file = e.target.files[0];
    if (file) {
      // const isValidSubtitle = fileName.endsWith('.ass') || fileName.endsWith('.srt');
      const isValidSubtitle = file.name.endsWith('.ass') || file.name.endsWith('.ASS');
      if (!isValidSubtitle) {
        setErrorMessage('请选择有效的字幕文件 (暂时只支持.ass 未来将支持.srt)');
        setTimeout(() => setErrorMessage(''), 3000);
        return;
      }
      setSubtitleFileName(file.name);
      // Process valid subtitle file
      const url = URL.createObjectURL(file);
      setSubtitleUrl(url);
    
      subParse(url, ".ass").then(parsed => {
        setErrorMessage(`已选择字幕文件: ${file.name}`);
        setTimeout(() => setErrorMessage(''), 3000);
        setAssLines(parsed);
      })

    }
  }

  const videoRef = useRef(null);
  const playerRef = useRef(null);
  const hideControlsTimeout = useRef(null);



  const togglePlay = useCallback(() => {
    console.log('Toggle play/pause');
    const video = videoRef.current;
    
    if (isPlaying) {
      video.pause();
    } else {
      video.play();
    }
    
    setIsPlaying(!isPlaying);
  }, [isPlaying]);
  useEffect(() => {
    if (curSubtitleLineIndex > 0) {
      renderJPSubtitles(curSubtitleLineIndex - 1, setPrevJPLine);
    }
    renderJPSubtitles(curSubtitleLineIndex, setCurJPLine);
    if (curSubtitleLineIndex < assLines.length - 1) {
      renderJPSubtitles(curSubtitleLineIndex + 1, setNextJPLine);
    }
    console.log(assLines);
    cache.current = {};
  }, [assLines]);
  const handleTimeUpdate = () => {
    setCurrentTime(videoRef.current.currentTime);
    let idx = rSubtitleLn.current;
    // console.log("Current time:", videoRef.current.currentTime);
    // console.log(`idx: ${idx} subtitle end time:`, assLines[idx]);
    while ( (idx < assLines.length - 1) && assLines[idx] && 
           (assLines[idx].end < +videoRef.current.currentTime + subtitleTimeOffset)) {
      // Update current subtitle line index
      // console.log(
      //   `Updating subtitle line index: 
      //   ${curSubtitleLineIndex} -> ${curSubtitleLineIndex + 1}, 
      //   current time: ${videoRef.current.currentTime}, 
      //   offset: ${subtitleTimeOffset}`
      // )
      idx++;
    }
    if (idx > rSubtitleLn.current) {
      rSubtitleLn.current = idx;
      if (idx > 0) {
        renderJPSubtitles(idx - 1, setPrevJPLine);
      }
      renderJPSubtitles(idx, setCurJPLine);
      if (idx < assLines.length - 1) {
        renderJPSubtitles(idx + 1, setNextJPLine);
      }
    }
    setCurSubtitleLineIndex(idx);

  };

  useEffect(() => {
    const intervalId = setInterval(() => {
      handleTimeUpdate();
    }, 100);
    return () => clearInterval(intervalId);
  }, [assLines])

  const handleDurationChange = () => {
    setDuration(videoRef.current.duration);
  };

  const handleSeek = (e) => {
    const seekTime = parseFloat(e.target.value);
    videoRef.current.currentTime = seekTime;
    setCurrentTime(seekTime);
  };

  const handleVolumeChange = (e) => {
    const newVolume = parseFloat(e.target.value);
    videoRef.current.volume = newVolume;
    setVolume(newVolume);
  };

  const handleFullscreen = () => {
    if (!document.fullscreenElement) {
      playerRef.current.requestFullscreen().catch(err => {
        setErrorMessage(`无法进入全屏模式: ${err.message}`);
      });
    } else {
      document.exitFullscreen();
    }
  };

  const [skipFeedback, setSkipFeedback] = useState({ show: false, direction: null });
  const skipFeedbackTimeout = useRef(null);
  // Function to show the skip feedback
  const showTimeSkipFeedback = (direction) => {
    setSkipFeedback({ show: true, direction });
    
    // Always show controls when skipping
    setShowControls(true);
    
    // Clear any existing timeout for the feedback
    if (skipFeedbackTimeout.current) {
      clearTimeout(skipFeedbackTimeout.current);
    }
    
    // Hide the feedback after 800ms
    skipFeedbackTimeout.current = setTimeout(() => {
      setSkipFeedback({ show: false, direction: null });
    }, 800);
  };

  const handleKeyDown = useCallback((e) => {
    // Only handle keyboard events if video is loaded
    if (!videoRef.current || !videoUrl) return;
  
    // Prevent default behavior for these keys (like scrolling with space)
    if (['Space', 'ArrowLeft', 'ArrowRight', ' '].includes(e.key)) {
      e.preventDefault();
    }
  
    switch (e.key) {
      case ' ':  // Space key
      case 'Space':
        togglePlay();
        break;
        
      case 'ArrowLeft':
        // Go back 5 seconds
        videoRef.current.currentTime = Math.max(0, videoRef.current.currentTime - 5);
        // console.log('Current time after rewind:', videoRef.current.currentTime);
        setCurrentTime(videoRef.current.currentTime);
        
        // Show visual feedback
        showTimeSkipFeedback('backward');
        break;
        
      case 'ArrowRight':
        // Go forward 5 seconds
        videoRef.current.currentTime = Math.min(
          duration, 
          videoRef.current.currentTime + 5
        );
        setCurrentTime(videoRef.current.currentTime);
        
        // Show visual feedback
        showTimeSkipFeedback('forward');
        break;
        
      default:
        break;
    }
  }, [videoUrl, duration, togglePlay]);

  const formatTime = (timeInSeconds) => {
    const hours = Math.floor(timeInSeconds / 3600);
    timeInSeconds %= 3600;
    const minutes = Math.floor(timeInSeconds / 60);
    const seconds = Math.floor(timeInSeconds % 60);
    if (hours > 0) {
      return `${hours}:${minutes < 10 ? '0' : ''}${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
    }
    return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`;
  };

  const handleVideoEnd = () => {
    setIsPlaying(false);
    setCurrentTime(0);
    videoRef.current.currentTime = 0;
  };

  const handlePlayerMouseMove = () => {
    setShowControls(true);
    
    // Clear any existing timeout
    if (hideControlsTimeout.current) {
      clearTimeout(hideControlsTimeout.current);
    }
    
    // Set new timeout to hide controls
    hideControlsTimeout.current = setTimeout(() => {
      if (isPlaying) {
        setShowControls(false);
      }
    }, 3000);
  };

  const handlePlayerMouseLeave = () => {
    if (isPlaying) {
      hideControlsTimeout.current = setTimeout(() => {
        setShowControls(false);
      }, 1000);
    }
  };

  useEffect(() => {
    document.addEventListener('fullscreenchange', () => {
      setFullscreen(!!document.fullscreenElement);
    });
    
    return () => {
      document.removeEventListener('fullscreenchange', () => {
        setFullscreen(false);
      });
      
      if (hideControlsTimeout.current) {
        clearTimeout(hideControlsTimeout.current);
      }
    };
  }, []);

  useEffect(() => {
    return () => {
      if (videoUrl && videoUrl.startsWith('blob:')) {
        URL.revokeObjectURL(videoUrl);
      }
    };
  }, [videoUrl]);

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    
    return () => {
      // Clean up event listeners
      window.removeEventListener('keydown', handleKeyDown);
      
      // Clean up timeouts
      if (hideControlsTimeout.current) {
        clearTimeout(hideControlsTimeout.current);
      }
      
      if (skipFeedbackTimeout.current) {
        clearTimeout(skipFeedbackTimeout.current);
      }
      
    };
  }, [duration, videoUrl, handleKeyDown]); 
  // window.addEventListener('keydown', handleKeyDown);

  // cached:
  const cache = useRef({});
  function renderJPSubtitles(idx, setter) {
    if (idx < 0 || idx >= assLines.length) return null;
    if (cache.current[idx]) {
      const cached = cache.current[idx];
      if (cached.setters.has(setter)) {
        return;
      }
      setter(cached.content);
      cached.setters.add(setter);
      return;
    } 
    cache.current[idx] = {
      content: null,
      setters: new Set([setter])
    };
    const cacheEntry = cache.current[idx];
    const line = assLines[idx];
    GetTokenizedString(line.jp_text).then(
      tokens => {
        let content;
        if (!tokens || tokens.length === 0) {
          console.log("Get tokens failed.")
          content = (<span>{line.jp_text || "null"}</span>);
        } else {
          console.log("Get tokens success:", tokens);
          content = (
            tokens.map((token, i) => (
              <HoverPreview
                key={`jp-token-${i}`}
                text={token}
                fGetContent={GetHTMLStringProxy}
              />
            ))
          )
        }
        cacheEntry.content = content;
        setter(content);

      }
    );
  }
  
  function renderCNSubtitles(idx) {
    if (idx < 0 || idx >= assLines.length) return null;
    const line = assLines[idx];
    return (
      <span>{ line.cn_text }</span>
    )
  }

  // useEffect(() => {
  //   if (curSubtitleLineIndex > 0) {
  //     renderJPSubtitles(curSubtitleLineIndex - 1, setPrevJPLine);
  //   }
  //   renderJPSubtitles(curSubtitleLineIndex, setCurJPLine);
  //   if (curSubtitleLineIndex < assLines.length - 1) {
  //     renderJPSubtitles(curSubtitleLineIndex + 1, setNextJPLine);
  //   }
  // }, [curSubtitleLineIndex, assLines]);

  return (
    <div className="video-player-page">
      <header className="video-header">
        <h1>视频播放器</h1>
        <h2>Video Player</h2>
      </header>
      {/* Add mode selector */}
      <div className="mode-selector">
        <button 
          className={`mode-button ${videoMode === VideoMode.WATCH ? 'active' : ''}`}
          onClick={() => setVideoMode(VideoMode.WATCH)}
        >
          观看模式
        </button>
        <button 
          className={`mode-button ${videoMode === VideoMode.STUDY ? 'active' : ''}`}
          onClick={() => setVideoMode(VideoMode.STUDY)}
        >
          学习模式
        </button>
        <div className="mode-help">
        <span className="help-icon material-icons">help_outline</span>
        <div className="tooltip">
          <div className="tooltip-content">
            <strong>观看模式:</strong> 大屏幕观看<br/>
            <strong>学习模式:</strong> 播放器缩小，实时AI讲解语法<br/>
          </div>
        </div>
      </div>
      </div>

      <div className="video-input-container">

        <div className="file-input-section">
          <p>选择本地视频文件：</p>
          <div className="custom-file-input">
            <input 
              type="file" 
              id="video-file" 
              accept="video/*" 
              onChange={handleFileSelect} 
              className="file-input" 
            />
            <label htmlFor="video-file"
              className={`file-input-label ${selectedFile ? 'has-file' : ''}`}>
              <span className="material-icons">upload_file</span>
              {selectedFile ? selectedFile.name : "选择文件"}
            </label>
          </div>
        </div>

        <div className="file-input-section">
          <p>选择本地字幕文件：</p>
          <div className="custom-file-input">
            <input
              type="file"
              id="subtitle-file"
              accept="*"
              onChange={handleSubtitleFileSelect}
              className="file-input"
            />
            <label htmlFor="subtitle-file"
              className={`file-input-label ${assLines.length > 0 ? 'has-file' : ''}`}>
              <span className="material-icons">upload_file</span>
              {assLines.length > 0 ? subtitleFileName : "选择字幕文件"}
            </label>
          </div>
        </div>

      </div>
      <div className='study-layout'>
      <div 
        className={`video-player-container ${videoMode} ${!showControls ? 'hide-controls' : ''}`}
        ref={playerRef}
        onMouseMove={handlePlayerMouseMove}
        onMouseLeave={handlePlayerMouseLeave}
        tabIndex="0"
        // onKeyDown={handleKeyDown}
      >
        <video
          ref={videoRef}
          src={videoUrl}
          onDurationChange={handleDurationChange}
          onEnded={handleVideoEnd}
          onClick={togglePlay}
          playsInline
        >
          您的浏览器不支持视频播放。
        </video>

        <div className={`video-controls ${!showControls ? 'hidden' : ''}`}>
          <div className="top-controls">
            <div className="time-display">
              {formatTime(currentTime)} / {formatTime(duration)}
            </div>
          </div>
          
          <div className="progress-container">
            <input
              type="range"
              className="progress"
              value={currentTime}
              max={duration || 0}
              onChange={handleSeek}
              step="0.1"
            />
          </div>
          
          <div className="bottom-controls">
            <button className="control-button" onClick={togglePlay}>
              {isPlaying ? (
                <span className="material-icons">pause</span>
              ) : (
                <span className="material-icons">play_arrow</span>
              )}
            </button>
            
            <div className="volume-container">
              <span className="material-icons">
                {volume === 0 ? 'volume_off' : 'volume_up'}
              </span>
              <input
                type="range"
                className="volume"
                min="0"
                max="1"
                step="0.1"
                value={volume}
                onChange={handleVolumeChange}
              />
            </div>
            
            <button className="control-button fullscreen" onClick={handleFullscreen}>
              <span className="material-icons">
                {fullscreen ? 'fullscreen_exit' : 'fullscreen'}
              </span>
            </button>
          </div>
        </div>

        {!videoUrl && (
          <div className="video-placeholder">
            <span className="material-icons">movie</span>
            <p>请选择本地文件播放</p>
          </div>
        )}
      </div>
      {videoMode === VideoMode.STUDY && (
        <div className='subtitle-display'>
          <div className='subtitle-block'>
            <div className='japanese-block'>
              {curSubtitleLineIndex > 0 && (
                <div className='subtitle prev'>
                  {prevJPLine}
                </div>
              )}
              <div className='subtitle current'>
                {curJPLine}
              </div>
              {curSubtitleLineIndex < assLines.length - 1 && (
                <div className='subtitle next'>
                  {nextJPLine}
                </div>
              )}
            </div>


          </div>
          <div className='subtitle-block'>
            <div className='chinese-block'>
              {curSubtitleLineIndex > 0 && (
                <div className='subtitle prev'>
                  {renderCNSubtitles(curSubtitleLineIndex - 1)}
                </div>
              )}
              <div className='subtitle current'>
                {renderCNSubtitles(curSubtitleLineIndex)}
              </div>
              {curSubtitleLineIndex < assLines.length - 1 && (
                <div className='subtitle next'>
                  {renderCNSubtitles(curSubtitleLineIndex + 1)}
                </div>
              )}
            </div>
          </div>
        </div>
        )}
        </div>



      <div className="video-info">
        <h3>使用说明:</h3>
        <ul>
          <li>选择一个本地视频播放</li>
          <li>按空格暂停/继续</li>
          <li>用左右箭头键快进/快退(5s)</li>
        </ul>
      </div>

      <div className="advance-settings">
        <h3>高级设置:</h3>
        <ul>
          <li>字幕时间轴基准: 
            <input
              type="number"
              style={{ width: '60px' }}
              value={0}
              onChange={(e) => {
                setSubtitleTimeOffset(parseFloat(e.target.value));
              }}
            ></input>

            (用于对齐视频与字幕，可以输入负值和小数，单位：秒)
          </li>
          <li>词典URL:
            <input
              type="text"
              value={dictUrl}
              disabled={true}
              onChange={(e) => setDictUrl(e.target.value)}
              placeholder="https://dict.youdao.com/result?word={}&lang=ja"
            ></input>
            (用于查询，使用{"{}"}作为单词的占位符。请注意：不同词典的HTML结构不同，不一定有效，最好通过插件形式实现其他词典的查询)


          </li>

        </ul>
        <button onClick={() => {
          GetTokenizedString("私は学生です").then(
            tokens => console.log(tokens)
          )
        }}>
          test
        </button>
      </div>
      {notification.show && (
          <Notification 
            message={notification.message} 
            type={notification.type} 
          />
      )}
      {/* test */}
      {/* <HoverPreview
        text="私"
        fGetContent={GetHTMLStringProxy}
      ></HoverPreview>
      <HoverPreview
        text="は"
        fGetContent={GetHTMLStringProxy}
      ></HoverPreview> */}
    </div>
  );
};

export default VideoPlayer;