import { h, render, FunctionalComponent, Fragment } from 'preact';
import { useState } from 'preact/hooks';

import { Button } from '@streatcodes/silk/components/button'
import { Input } from '@streatcodes/silk/components/input'
import './login.css';

import * as api from './api';

export const Login: FunctionalComponent = () => {
    const [loading, setLoading] = useState(false);
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const onLogin = async (e: Event) => {
        e.preventDefault();
        setLoading(true);

        try {
            const token = await api.login(email, password);
            console.log(token)
            console.log('succesful login')
        } catch (e) {
            //TODO show error
            console.error(e.message)
        }

        setLoading(false);
    }

    return (
        <div className="login">
            <form onSubmit={onLogin}>
                <label>
                    <p>Email</p>
                    <Input type="email" value={email} onChange={v => setEmail(v as string)} />
                </label>
                <label>
                    <p>Password</p>
                    <Input type="password" value={password} onChange={v => setPassword(v as string)} />
                </label>
                <Button loading={loading}>Login</Button>
            </form>
        </div>
    );
};