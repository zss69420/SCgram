import React, { useState } from "react";
import { Route, Switch, Redirect } from "react-router-dom";

import Login from "./Login";
import Register from "./Register";
import Home from "./Home";

//router:
// path: /login  -> component: Login
//       /register -> component: Register
//       /home ->  component: Home
function Main(props) {
    const { isLoggedIn, handleLoggedIn } = props;

    //case 1: logged in -> display home page
    //cas2 2: otherwise, display login page
    const showLogin = () => {
        return isLoggedIn ? (
            //redirect to home page
            <Redirect to="/home" />
        ) : (
            <Login handleLoggedIn={handleLoggedIn} />
        );
    };

    //showHome is also affected by isLoggedIn because when it is not logged in
    //we cannot show home directly, instead we have to redirect to login page letting user to login first
    const showHome = () => {
        //case 1: logged in, display home page
        //case 2: otherwise, display login page
        return isLoggedIn ? <Home /> : <Redirect to="/login" />;
    };

    //switch is to just keep only one route can be matched at one time
    return (
        <div className="main">
            <Switch>
                <Route path="/" exact render={showLogin} />
                <Route path="/login" render={showLogin} />
                <Route path="/register" component={Register} />
                <Route path="/home" render={showHome} />
            </Switch>
        </div>
    );
}

export default Main;
