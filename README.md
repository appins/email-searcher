# email-searcher
Given a domain name and real name, attempt to find an existing email for that user.

### Using
Run it with both the domain and name flags, like:
```
./email-searcher --name='alex anderson' --domain='appins.dev'
```

It will attempt to guess the email and report back any successful findings. 
For example, the above query will yield my correct email: alex@appins.dev,
even if it is only able to determine that its probably the correct email.
