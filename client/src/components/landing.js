import React, { Component, useState } from 'react';

function Landing(props) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loginErrorMessage, setLoginErrorMessage] = useState('');
  const [registerErrorMessage, setRegisterErrorMessage] = useState('');


  const [tab, setTab] = useState('about');

  const handleSubmit = () => {
    const data = new FormData();
    data.append('username', username);
    data.append('password', password);

    if (tab === 'login') {
      fetch('/login', {
        body: data,
        method: 'post'
      }).then(res => {
        if (res.status === 200) {
          // redirect to game
          window.location.href = '/game';
        } else if (res.status === 400) {
          setLoginErrorMessage('There was a problem logging in. Please use a correct username and password.');
        }
      });
    } else if (tab === 'register') {
      fetch('/register', {
        body: data,
        method: 'post'
      }).then(res => {
        if (res.status === 200) {
          window.location.href = '/game';
        } else {
          setRegisterErrorMessage('There was a problem registering. Please pick a novel username.');
        }
      });
    }
  };

  return (
    <div>
      <h1>Dungeon Online</h1>

      <div id="wrapper">
        <div className="box">
          <span className="tab" onClick={() => setTab('about')}>about</span> |{' '}
          <span className="tab" onClick={() => setTab('register')}>register</span> |{' '}
          <span className="tab" onClick={() => setTab('login')}>login</span>
        </div>

        <div className="box-no-border">
          {tab === 'about' &&
            <div>
              Play my cool ass game dude. It's cool.
            </div>
          }

          {tab === 'register' && <h2>Register</h2>}

          {tab === 'login' && <h2>Login</h2>}

          {(tab === 'register' || tab === 'login') &&
              <>
                <input type="text" placeholder="Username" onChange={e => setUsername(e.target.value)} value={username} />
                <input type="password" placeholder="Password" onChange={e => setPassword(e.target.value)} value={password} />
              </>
          }

          {tab === 'register' && <input type="submit" value="Register" onClick={handleSubmit} />}

          {tab === 'login' && <input type="submit" value="Login" onClick={handleSubmit} />}
        
          {(tab === 'login' && loginErrorMessage) && <div className="box">{loginErrorMessage}</div>}
          {(tab === 'register' && registerErrorMessage) && <div className="box">{registerErrorMessage}</div>}
        </div>
      </div>
    </div>
  );
}

export default Landing;