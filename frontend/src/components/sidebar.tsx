import { h, render, FunctionalComponent, Fragment } from 'preact';
import { Link } from 'preact-router/match';
import './sidebar.css';

export const SideBar: FunctionalComponent = () => {

    return <div className="side-bar">
        <Link href="/" activeClassName="active" className="link">Overview</Link>
        <Link href="/hosts" activeClassName="active" className="link">Hosts</Link>
        <Link href="/metrics" activeClassName="active" className="link">Metrics</Link>
    </div>
}