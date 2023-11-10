import SlackButton from "./SlackButton"
import Header from "./Header"
import Footer from "./Footer"
import Head from "./Head"

type FAQProps = {
  success: boolean
}

export default function FAQ(props: FAQProps) {
  return (
    <div>
      <Head>
        <title>FAQ – Standup Jack</title>
      </Head>
      <Header {...props} />
      <div className="faq">
        <div class="question">
          <h2>How much does Standup Jack cost?</h2>
          <p>
            Jack is free for everyone to try for the first 14 days. After that
            Jack costs <strong>$1.00/user/month</strong>.
          </p>
          <p>
            Jack will message you when you're getting close to the end of your
            trial period with details on how to upgrade.
          </p>
          <p>
            If you choose to let your trial expire, Jack will gracefully go into
            hibernation until you choose to upgrade or decide to remove him.
          </p>
        </div>

        <div class="question">
          <h2>Does Standup Jack support different timezones?</h2>
          <p>
            Yes of course! Jack looks up the timezone setting in your Slack
            profile and makes sure that 10am is not 3am where you live.
          </p>
          <p>
            You'll definitely want to double-check that the timezone settings
            are up-to-date, because Slack doesn't do this automatically for you.
          </p>
        </div>

        <div class="question">
          <h2>Can I change the time I get asked for an update?</h2>
          <p>
            Of course! You just need to say <strong>change my time</strong>,
            then pick your time and confirm.
          </p>
        </div>

        <div class="question">
          <h2>Can I change the time when standup gets posted?</h2>
          <p>
            Certainly! You just need to say <strong>change standup time</strong>
            , then pick a new time.
          </p>
        </div>

        <div class="question">
          <h2>Can I have standups on different days?</h2>
          <p>
            Definitely! You just need to say <strong>change my time</strong>,
            then include the days you want to have standup. For example,{" "}
            <strong>10am mondays, wednesdays and fridays</strong>
          </p>
        </div>

        <div class="question">
          <h2>Can I change who is participating in standup?</h2>
          <p>
            Absolutely! You just need to say <strong>add a user</strong> to add
            your teammates or <strong>remove a user</strong> to remove a
            teammate.
          </p>
        </div>

        <div class="question">
          <h2>Can I post team standup to a different channel?</h2>
          <p>
            Yep! You just need to say <strong>change standup channel</strong>,
            then pick which channel you want to post to.
          </p>
        </div>

        <div class="question">
          <h2>Can I have multiple standups per team?</h2>
          <p>
            Absolutely! You just need to say <strong>create a standup</strong>{" "}
            to create a new standup
          </p>
        </div>

        <div class="question">
          <h2>Can I pick my own questions?</h2>
          <p>
            Most def! You can add questions, change questions or remove
            questions. To pick your own questions, message Jack:
          </p>
          <ul>
            <li>
              <strong>add question</strong> to add a new question to standup
            </li>
            <li>
              <strong>change question</strong> to change an existing question
            </li>
            <li>
              <strong>delete question</strong> to remove a question from standup
            </li>
          </ul>
        </div>

        <div class="question">
          <h2>Can I tell Jack to post into private groups?</h2>
          <p>
            Yep! Whenever you create a new standup, you'll have the option to
            invite @jack to a private group. There's a couple more steps, but
            Jack will walk you through what you need to do to get that setup.
          </p>
        </div>

        <div class="question">
          <h2>Where can I read your privacy policy?</h2>
          <p>
            You can visit <a href="/privacy">our privacy policy</a> to learn
            more about what data we retain, for how long and how you can request
            data removal. interested! Right now Jack is only available
            <p>
              If you have any questions or concerns, please don't hesitate to{" "}
              <a href={"mailto:hi@standupjack.com"}>reach out</a>!
            </p>
          </p>
        </div>

        <div class="question">
          <h2>I love this idea! But our team doesn't use Slack.</h2>
          <p>
            Happy to hear you're interested! Right now Jack is only available on
            Slack, but <a href={"mailto:hi@standupjack.com"}>contact me</a> if
            you'd like Jack on your system. I will build additional adapters for
            email, hipchat, IRC or any other system depending on demand.
          </p>
        </div>

        <div class="question">
          <h2>What type of animal is Jack?</h2>
          <p>
            Jack is a rooster. But he doesn't consider himself to be like the
            other roosters. His friends say he's a good listener and very
            respectful.
          </p>
        </div>
        <div class="buttons">
          <span class="subtext">Ready to communicate better?</span>
          <SlackButton jingle={true} {...props} />
        </div>
      </div>
      <Footer {...props} />
    </div>
  )
}
