## Go-Home

A TUI that believes that the only thing that matters is the next 7 days. Go-Home was created and designed
with the idea that most of the time we don't need to know what happens in the next 2 months or 3 weeks, we
just need to know what I need to do today and a bit into the future.

Using the power of Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) we have in our hands
a simple way to check up on, lightly edit, and delete events all within the wonderful world of the terminal.

## Setup

0. clone the repo, build, then run it once to populate config files

1. Go to **[GCP](https://console.cloud.google.com/)** logging in to the account which is associated with the google calendar you would like to use.

2. Create a new project under any name you like, but to keep it organized I recommend "go-home"

3. Under "Quick access" in the home page or via the navigation menu on the top left or opened with ".", then go to **"APIs & Services"**

   <img width="1108" height="249" alt="image" src="https://github.com/user-attachments/assets/f4add942-1ea2-41dd-aaf6-11e8e8f1b95e" />
   <img width="544" height="860" alt="image" src="https://github.com/user-attachments/assets/bbe1069d-b70b-4e2f-abba-1e31db9d3421" />

4. Select **" + Enable APIs and services"**

   <img width="888" height="159" alt="image" src="https://github.com/user-attachments/assets/c54bd864-be72-47e1-b8a8-448172c08161" />

5. Scroll down and select **"Google Calendar API"** or type it in the search bar

6. Enable the Google Calendar API service

   <img width="484" height="210" alt="image" src="https://github.com/user-attachments/assets/df5e9d95-2b41-437f-8f95-723f4afe4273" />

7. On the left sidebar select **"Credentials"** still within the APIs & Services.

   <img width="670" height="270" alt="image" src="https://github.com/user-attachments/assets/a74f073d-6b2e-468b-86bb-9f4ab8571d00" />

8. Under **"+ Create credentials"** select **"OAuth client ID"**

   <img width="639" height="283" alt="image" src="https://github.com/user-attachments/assets/48162d7c-7fa2-48d1-bc9b-799e4d2d033f" />

9. Create a new OAuth client ID with an application type of **"Desktop app"** and any name you prefer for the client

   <img width="549" height="402" alt="image" src="https://github.com/user-attachments/assets/8a9c8663-9eaa-4406-9f59-4f236f0c5d52" />

10. When the popup shows you can manually copy both **"Client ID"** and **"Client secret"** or download the JSON file. We will need this later

    <img width="543" height="673" alt="image" src="https://github.com/user-attachments/assets/8ad16955-9479-477c-90de-b6223525837b" />

11. Select your newly created OAuth Client then go to **"Branding"** on the lefthand side

    <img width="1128" height="144" alt="image" src="https://github.com/user-attachments/assets/2c928336-b1bd-4e66-9d7a-b7c4a79468d0" />
    <img width="283" height="376" alt="image" src="https://github.com/user-attachments/assets/d5f0061d-7fc9-4af8-87a0-1a91ea4fc2d5" />

12. Update fields **"User support email"** and **"Email address"** under Developer contact information. You can simply put the email associated with the google calendar, then save changes. This is for the OAuth web page that opens doing authentication. We do this since GCP has the impression this will be a service for other users but it's just for us, so we must comply

13. Select **"Audience"** on the lefthand side

    <img width="252" height="324" alt="image" src="https://github.com/user-attachments/assets/c2cbd513-925b-4c3b-956b-8d0e3a49a893" />

14. Add a test user under the same gmail account as previously used and with that we are done with GCP

15. Going to you config location

- Windows - %AppData%/go-home/.env
- Mac - /Library/Application Support/go-home/.env
- Linux - /.config/go-home/.env

and populate the following

- CALENDAR_ID="myemail@gmail.com"
- CLIENT_ID="from step 10"
- CLIENT_SECRET="from step 10"

16. Lastly using the flag -a (auth) go through google authentication using the same email as before. Do note
    it will say the application is not verified, this is the byproduct of again Google assuming this is a large
    application for many users and we don't really care if it's verified because it's for us
