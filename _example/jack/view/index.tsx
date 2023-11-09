import Header from "./Header"
import SlackButton from "./SlackButton"
import Head from "./Head"
import Footer from "./Footer"

type IndexProps = {
  url?: {
    query: {
      success?: boolean
    }
  }
}

export default function Index(props: IndexProps = {}, context = {}) {
  const url = props.url || {
    query: {},
  }
  return (
    <div className="index page">
      <Head>
        <title>Standup Jack</title>
      </Head>
      <Header success={url.query.success} />
      <Walkthrough success={url.query.success} />
      <div className="footer">
        <Footer success={url.query.success} />
      </div>
    </div>
  )
}

type WalkthroughProps = {
  success?: boolean
}

function Walkthrough({ success }: WalkthroughProps) {
  return (
    <div className="Walkthrough">
      <ul className="list">
        <li className="item">
          <div className="count">1</div>
          <div className="left">
            <h3>Add Jack to Slack</h3>
            <p>This will only take a minute.</p>
            <p>
              Jack helps you plan your day and keep your teammates in the loop.
            </p>
          </div>
          <div className="right">
            <div className="button-container">
              <SlackButton success={success} jingle={true} />
            </div>
            <div class="feature-list">
              <div className="feature-inner">
                <div class="feature">
                  <img src="/images/checkmark-green.svg" />5 minute setup
                </div>
                <div class="feature">
                  <img src="/images/checkmark-green.svg" />
                  Remote ready
                </div>
                <div class="feature">
                  <img src="/images/checkmark-green.svg" />
                  Manage inside Slack
                </div>
                <div class="feature">
                  <img src="/images/checkmark-green.svg" />
                  14 days for free
                </div>
                <div class="feature">
                  <img src="/images/checkmark-green.svg" />
                  $1 per user per month
                </div>
              </div>
            </div>
          </div>
        </li>
        <li className="item">
          <div className="count">2</div>
          <div className="left">
            <h3>Create a Standup</h3>
            <p>
              Once you add Jack to Slack, he will ask you what time you want to
              have standup, with whom, and where to post the updates.
            </p>
            <p>
              Jack will then introduce himself to each teammate and ask them to
              pick a time that works best for them to give their update.
            </p>
          </div>
          <div className="right">
            <div className="image-container">
              <img src="/images/create-standup.png" alt="create standup" />
            </div>
          </div>
        </li>
        <li className="item">
          <div className="count">3</div>
          <div className="left">
            <h3>Give Jack your Daily Update</h3>
            <p>
              Jack will message you at the time you picked and ask you a few
              questions.
            </p>
            <p>
              These questions are designed to help you stay on track and let
              your teammates know what you're up to.
            </p>
          </div>
          <div className="right">
            <div className="image-container">
              <img src="/images/give-update.png" alt="give update" />
            </div>
          </div>
        </li>
        <li className="item">
          <div className="count">4</div>
          <div className="left">
            <h3>Get your Team's Updates</h3>
            <p>
              Every weekday at the time you picked, Jack will post everyone's
              responses to the channel you chose.
            </p>
          </div>
          <div className="right">
            <div className="image-container">
              <img src="/images/post-standup.png" alt="post standup" />
            </div>
          </div>
        </li>
        <li className="item" data-final>
          <div className="count checkmark" />
          <div className="left">
            <h3>You're All Set!</h3>
            <p className="casablanca">
              <a
                href="https://www.youtube.com/watch?v=DDybg9CNXcM"
                target="_blank"
              >
                I think this is the beginning of a beautiful friendship.
              </a>
            </p>
          </div>
          <div className="right">
            <div class="slack-button">
              <SlackButton success={success} jingle={true} />
              <div class="privacy">
                <a href="/privacy">Privacy Policy</a>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  )
}
