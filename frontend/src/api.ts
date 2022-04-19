export async function login(email: string, password: string): Promise<string> {
    const req = {
        email,
        password
    }
    
    const res = await fetch(`/api/login`, {
        method: 'POST',
        body: JSON.stringify(req)
    });

    if(res.status > 299) {
        throw new Error(`Unexpected status code ${res.status}`)
    }

    const data = await res.json();

    return data.token;
}