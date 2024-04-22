# Stori Card Challenge

Hi, I prepared this challenge for the guys from StoriCard

## Problem and proposed solution

StoriCard needs a way to generate account resumes from users, these resumes will be sent by mail and should share information about the last few months' transactions and account balance.

For this I will create a worker in golang that will read a CSV file, this worker will store the data in a SQLite database and send an email with the resume.


## The Mail

### Desktop

![imagen](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/940e94bd-817b-40a5-b72c-5f27da977184)

### Mobile

![WhatsApp Image 2024-04-22 at 12 42 49 AM](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/5228f1ae-0a58-42a9-a965-48ee14c97f66)

## How I imagine this on producction

The solution I thought of is quite basic because it is only to solve the specific exercise.

Despite that, I took the time to prepare, as I imagine, a similar solution, for a productive environment so that this development can scale correctly.

![imagen](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/553f025c-0070-4b28-80b4-e700362f0f26)

## Extra feature: Whatsapp

I added a module to send Message using the public whatsapp business API

![WhatsApp Image 2024-04-22 at 1 55 14 AM](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/3737f387-288f-437a-834b-5fa0a38a27db)

## Sketches

Some ugly sketches from the process

![imagen](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/39f60310-4030-4975-a854-4296117e60e4)
![WhatsApp Image 2024-04-22 at 12 08 58 AM](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/6d72f000-45aa-4e6a-86c2-e714987a56ec)
![WhatsApp Image 2024-04-22 at 12 08 58 AM(1)](https://github.com/LucasRosello/stori-card-challenge/assets/55340118/034cf565-d983-4147-8035-9e3fc1b87ea6)
