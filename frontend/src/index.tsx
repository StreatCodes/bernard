import { h, render, FunctionalComponent, Fragment } from 'preact';

import { Button } from '@streatcodes/silk/components/button'
import { Login } from './login';

const App: FunctionalComponent = () => {

    return (
        <main>
            <Login />
        </main>
    );
};

const appElement = document.getElementById('app');
if (appElement) {
    render(<App />, appElement);
} else {
    throw new Error(`No element found with id 'app'`);
}