import React, { Component, useEffect, useState } from 'react';
import {keyDown, keyUp} from '../redux/actions';

function Input(props) {

  useEffect(() => {
    const handleKeyDown = event => {
      props.dispatch(keyDown(event.code));
    }
  
    const handleKeyUp = event => {
      props.dispatch(keyUp(event.code));
    }

    window.addEventListener('keydown', handleKeyDown);
    window.addEventListener('keyup', handleKeyUp);
  }, []);

  return null;
}

export default Input;