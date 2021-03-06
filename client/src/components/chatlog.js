import React, { Component, useRef, useEffect, useState } from 'react';
import {sendChat, setTyping} from '../redux/actions';

function ChatLog(props) {
  const {log: {log}} = props;

  const messagesEndRef = useRef(null);

  useEffect(() => {
    messagesEndRef.current.scrollIntoView();
  }, [JSON.stringify(log)]); // seriously???

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
          {log.map((entry, i) => {
              if (entry.serverMessage) {
                return <li key={i} style={{color: '#adadad'}}>{entry.serverMessage.message}</li>
              }

              if (entry.chat) {
                return <li key={i}><span style={{color: '#00ff00'}}>{entry.chat.name}</span>: {entry.chat.message}</li>;
              }

              if (entry.spawn) {
                return <li key={i} style={{color: '#adadad'}}><span style={{color: '#00ff00'}}>{entry.spawn.name}</span> entered the zone.</li>;
              }

              if (entry.despawn) {
                return <li key={i} style={{color: '#adadad'}}><span style={{color: '#00ff00'}}>{entry.despawn.name}</span> left the zone.</li>;
              }

              if (entry.attack) {
                if (entry.attack.hit) {
                  return <li key={i} style={{color: '#adadad'}}><span style={{color: '#00ff00'}}>{entry.attack.attacker}</span> hit <span style={{color: '#00ff00'}}>{entry.attack.target}</span> for <span style={{color: '#ff0000'}}>{entry.attack.damage}</span> damage!</li>
                }
                return <li key={i} style={{color: '#adadad'}}><span style={{color: '#00ff00'}}>{entry.attack.attacker}</span> missed <span style={{color: '#00ff00'}}>{entry.attack.target}</span>!</li>
              }

              if (entry.die) {
                return <li key={i} style={{color: '#adadad'}}><span style={{color: '#00ff00'}}>{entry.die.name}</span> was slain.</li>
              }

              if (entry.zone) {
                return <li key={i} style={{color: '#adadad'}}>You entered {entry.zone.name}.</li>
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