import React, { Component, useRef, useEffect, useState } from 'react';
import {sendChat, setTyping} from '../redux/actions';

function ChatLog(props) {
  const {messages} = props;

  const messagesEndRef = useRef(null);

  useEffect(() => {
    messagesEndRef.current.scrollIntoView();
  }, [JSON.stringify(messages)]); // seriously???

  const [chatInput, setChatInput] = useState('');

  const handleKeyDown = (event) => {
    if (event.key === 'Enter') {
      props.dispatch(sendChat(chatInput));
      setChatInput('');
    }
  }

  const handleFocus = (event) => {
    props.dispatch(setTyping(true));
  }

  const handleFocusOut = (event) => {
    props.dispatch(setTyping(false));
  }

  return (
    <div id="chat-log-wrapper">
      <div id="chat-log">
        <ul>
          {messages.map((message, i) => {
              if (message.chat) {
                return <li key={i}><span style={{color: '#00ff00'}}>{message.chat.From}</span>: {message.chat.Message}</li>;
              }

              if (message.join) {
                return <li key={i} style={{color: '#ff0000'}}>{message.join.From} joined.</li>;
              }
          })}
          <div ref={messagesEndRef} />
        </ul>
        
        <input type="text" onChange={e => setChatInput(e.target.value)} onKeyDown={handleKeyDown} onFocus={handleFocus} onBlur={handleFocusOut} value={chatInput}  />
      </div>
    </div>
  );
}

export default ChatLog;