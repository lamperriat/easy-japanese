import React, { useState, useRef, useEffect, useCallback } from 'react';
import './VideoPlayer.css';

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

  const handleTimeUpdate = () => {
    setCurrentTime(videoRef.current.currentTime);
  };

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
    const minutes = Math.floor(timeInSeconds / 60);
    const seconds = Math.floor(timeInSeconds % 60);
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
          <p>选择本地文件：</p>
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
        
        {errorMessage && (
          <div className={`message ${errorMessage.includes('已选择') ? 'success-message' : 'error-message'}`}>
            {errorMessage}
          </div>
        )}
      </div>

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
          onTimeUpdate={handleTimeUpdate}
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

      <div className="video-info">
        <h3>使用说明:</h3>
        <ul>
          <li>选择一个本地视频播放</li>
          <li>按空格暂停/继续</li>
          <li>用左右箭头键快进/快退(5s)</li>
        </ul>
      </div>
    </div>
  );
};

export default VideoPlayer;