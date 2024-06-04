import { common } from '@mui/material/colors';
import { stringify } from 'querystring';
import React, { useState ,useEffect } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import * as myCommon from '../common.js';

function LoginPage() {
    const [email, setEmail] = useState('');
    const [redirectToPassword, setRedirectToPassword] = useState(false);
    const [socialIdps, setSocialIdps] = useState<{ slug: string }[]>([]);

    useEffect(() => {
      const fetchData = async () => {
        const response = await fetch("http://localhost1.com:9044/api/manifest", {
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
  
    const handleSubmit = async (event: React.FormEvent) => {
        event.preventDefault();
        let csrf = myCommon.getCSRF();
        console.log(csrf);
        // Fetch call to validate the email
        const response = await fetch("http://localhost1.com:9044/api/login-phase-one", {
          method: "POST",
          credentials: 'include', 
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
            "X-Csrf-Token":csrf,

          },
          body: JSON.stringify({ 
            email: email,
             }), // send the email as part of the request body
        });
      
        if (response.ok) {
          const data = await response.json();
          if (data.isValid) {
            setRedirectToPassword(true);
          } else {
            // Handle invalid email
          }
        } else {
          // Handle error
        }
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