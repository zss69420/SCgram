import logo from '../assets/images/logo.svg';
import '../styles/App.css';
import TopBar from "./TopBar";
import React, { useState } from "react";
import { TOKEN_KEY } from "../constants";
import Main from "./Main";

function App() {
  //after getting token from local storage, and destructure it
  const [isLoggedIn, setIsLoggedIn] = useState(
      localStorage.getItem(TOKEN_KEY) ? true : false
  );

  //define this logout func, and remove token, and set logged to false;
  //pass this information to child component: TopBar, because we have a logged out button in TopBar
  //when you click it, we have remove token and set un-logged-in
  const logout = () => {
    console.log("log out");
    localStorage.removeItem(TOKEN_KEY);
    setIsLoggedIn(false);
  };

  //define this callback func loggedIn which can pass to Main, then to Login component
  //then Login can call this cb and bring token all the way back to App
  //and then set token: {key = TOKEN_KEY: value = token}
  const loggedIn = (token) => {
    if (token) {
      localStorage.setItem(TOKEN_KEY, token);
      setIsLoggedIn(true);
    }
  };

  return (
      <div className="App">
        <TopBar isLoggedIn={isLoggedIn} handleLogout={logout} />
        <Main isLoggedIn={isLoggedIn} handleLoggedIn={loggedIn} />
      </div>
  );
}

export default App;





