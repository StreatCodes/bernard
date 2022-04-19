import { h, render, FunctionalComponent, Fragment } from 'preact';
import { Router, route } from 'preact-router';

import './index.css';

import { Button } from '@streatcodes/silk/components/button'
import { Login } from './routes/login';
import { useEffect, useState } from 'preact/hooks';
import { Overview } from './routes/overview';
import { Hosts } from './routes/hosts';
import { Metrics } from './routes/metrics';
import { SideBar } from './components/sidebar';

const App: FunctionalComponent = () => {
    const [isAuthed, setIsAuthed] = useState(true);

    useEffect(() => {
        const token = localStorage.getItem('bernard-token');
        setIsAuthed(token !== null)
    }, [])

    if (!isAuthed) return <Login isAuthed={setIsAuthed} />

    return (
        <main>
            <SideBar />
            <Router>
                <Overview path="/" />
                <Hosts path="/hosts" />
                <Metrics path="/metrics" />
            </Router>
        </main>
    );
};

const appElement = document.getElementById('app');
if (appElement) {
    render(<App />, appElement);
} else {
    throw new Error(`No element found with id 'app'`);
}

