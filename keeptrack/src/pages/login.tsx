import { stringify } from 'querystring';
import React, { useState ,useEffect } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
function LoginPage() {
    const [email, setEmail] = useState('');
    const [redirectToPassword, setRedirectToPassword] = useState(false);
    const [socialIdps, setSocialIdps] = useState<{ slug: string }[]>([]);

    useEffect(() => {
      const fetchData = async () => {
        const response = await fetch("http://localhost:9044/api/manifest", {
          method: "GET",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
          },
        });
  
        if (response.ok) {
          const data = await response.json();
          setSocialIdps(data.social_idps);
        } else {
          // Handle error
        }
      };
  
      fetchData();
    }, []);
  
    const handleSubmit = (event: React.FormEvent) => {
      event.preventDefault();
      setRedirectToPassword(true);
    };
  
    return (
      <form onSubmit={handleSubmit}>
        <h1>Login</h1>
        <label>
          Email:
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </label>
        <button type="submit">Next</button>
        <div>
        <a href="/signup">Sign Up</a> | <a href="/forgot-password">Forgot Password?</a>
      </div>
        <div style={{ display: 'flex', flexDirection: 'row', justifyContent: 'space-between' }}>
        {socialIdps.map((idp) => (
          <button key={idp.slug}>{idp.slug}</button>
        ))}
        </div>
      </form>
    );
  }
  

export default LoginPage;