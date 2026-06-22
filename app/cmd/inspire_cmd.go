package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var InspireCommand = &cobra.Command{
	Use:   "inspire",
	Short: "Return an Islamic inspirational quote",
	Run:   handleInspireCmd(),
}

func handleInspireCmd() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		quotes := []string{
			"Indeed, with hardship [will be] ease. – [Qur'an 94:6]",
			"So remember Me; I will remember you. – [Qur'an 2:152]",
			"And He found you lost and guided [you]. – [Qur'an 93:7]",
			"Verily, in the remembrance of Allah do hearts find rest. – [Qur'an 13:28]",
			"And whoever puts his trust in Allah, then He will suffice him. – [Qur'an 65:3]",
			"Allah does not burden a soul beyond that it can bear. – [Qur'an 2:286]",
			"So be patient. Indeed, the promise of Allah is truth. – [Qur'an 30:60]",
			"And rely upon Allah; and sufficient is Allah as Disposer of affairs. – [Qur'an 33:3]",
			"And whoever fears Allah – He will make for him a way out. – [Qur'an 65:2]",
			"O you who have believed, seek help through patience and prayer. – [Qur'an 2:153]",
			"Whoever does righteousness, male or female, while he is a believer – We will surely cause him to live a good life. – [Qur'an 16:97]",
			"And be not like those who forgot Allah, so He made them forget themselves. – [Qur'an 59:19]",
			"Say, ‘My Lord has commanded justice.’ – [Qur'an 7:29]",
			"Indeed, the mercy of Allah is near to the doers of good. – [Qur'an 7:56]",
			"And your Lord is going to give you, and you will be satisfied. – [Qur'an 93:5]",
			"Say, 'Nothing will happen to us except what Allah has decreed for us: He is our protector.' – [Qur'an 9:51]",
			"The best among you are those who have the best manners and character. – [Hadith, Bukhari]",
			"The strong person is not the one who can overpower others, but the one who controls himself when angry. – [Hadith, Bukhari]",
			"Allah is kind and loves kindness in all matters. – [Hadith, Bukhari & Muslim]",
			"Make things easy for people and do not make them difficult, and cheer them up and do not drive them away. – [Hadith, Bukhari]",
			"Indeed, Allah loves those who rely upon Him. – [Qur'an 3:159]",
			"My mercy encompasses all things. – [Qur'an 7:156]",
			"And speak to people good [words]. – [Qur'an 2:83]",
			"And We have not sent you, [O Muhammad], except as a mercy to the worlds. – [Qur'an 21:107]",
			"So do not weaken and do not grieve, and you will be superior if you are [true] believers. – [Qur'an 3:139]",
			"Indeed, prayer prohibits immorality and wrongdoing. – [Qur'an 29:45]",
			"Whoever comes [on the Day of Judgment] with a good deed will have ten times the like thereof [to his credit]. – [Qur'an 6:160]",
			"And He is with you wherever you are. – [Qur'an 57:4]",
			"Your ally is none but Allah and [therefore] His Messenger and those who have believed. – [Qur'an 5:55]",
			"Every soul will taste death. – [Qur'an 3:185]",
			"Do not lose hope in the mercy of Allah. – [Qur'an 39:53]",
			"He created death and life to test you as to which of you is best in deed. – [Qur'an 67:2]",
			"Call upon Me; I will respond to you. – [Qur'an 40:60]",
			"Indeed, Allah is with the patient. – [Qur'an 2:153]",
			"The most beloved deeds to Allah are those done consistently, even if small. – [Hadith, Bukhari & Muslim]",
			"None of you truly believes until he loves for his brother what he loves for himself. – [Hadith, Bukhari & Muslim]",
			"The world is a prison for the believer and a paradise for the disbeliever. – [Hadith, Muslim]",
			"Whoever believes in Allah and the Last Day should speak good or remain silent. – [Hadith, Bukhari & Muslim]",
			"Part of someone's being a good Muslim is leaving alone that which does not concern him. – [Hadith, Tirmidhi]",
			"Be mindful of Allah, and Allah will protect you. – [Hadith, Tirmidhi]",
			"Guidance is not but from Allah. – [Qur'an 28:56]",
			"And establish prayer and give zakah and obey the Messenger – that you may receive mercy. – [Qur'an 24:56]",
			"And whoever turns away from My remembrance – indeed, he will have a depressed life. – [Qur'an 20:124]",
			"Indeed, Allah does not wrong the people at all, but it is the people who are wronging themselves. – [Qur'an 10:44]",
			"And whoever is grateful – his gratitude is only for [the benefit of] himself. – [Qur'an 31:12]",
			"Indeed, Allah commands justice and good conduct and giving to relatives. – [Qur'an 16:90]",
			"Indeed, Allah is ever Knowing and Wise. – [Qur'an 4:17]",
			"And We will surely test you with something of fear and hunger… but give good tidings to the patient. – [Qur'an 2:155]",
			"And your Lord is most forgiving, owner of mercy. – [Qur'an 18:58]",
			"Say: 'In the bounty of Allah and in His mercy – in that let them rejoice.' – [Qur'an 10:58]",
			"Indeed, Allah is ever Acquainted and Seeing. – [Qur'an 4:35]",
			"Truly, no one despairs of Allah’s soothing mercy except those who have no faith. – [Qur'an 12:87]",
			"Whoever humbles himself for the sake of Allah, Allah will raise him. – [Hadith, Muslim]",
			"Visit the sick, feed the hungry, and free the captive. – [Hadith, Bukhari]",
			"Smiling in the face of your brother is charity. – [Hadith, Tirmidhi]",
			"Allah does not look at your appearance or wealth but at your hearts and deeds. – [Hadith, Muslim]",
			"The best of you are those who are best to their families. – [Hadith, Tirmidhi]",
			"He who does not show mercy will not be shown mercy. – [Hadith, Bukhari & Muslim]",
			"Make repentance, for I make repentance a hundred times a day. – [Hadith, Muslim]",
			"Charity does not decrease wealth. – [Hadith, Muslim]",
			"Whoever remembers Allah, Allah remembers him. – [Qur'an 2:152]",
			"And be moderate in your pace and lower your voice; indeed, the most disagreeable of sounds is the voice of donkeys. – [Qur'an 31:19]",
			"Indeed, good deeds erase bad deeds. – [Qur'an 11:114]",
			"Do not walk on the earth arrogantly. – [Qur'an 17:37]",
			"Indeed, those who have said, 'Our Lord is Allah' and then remained steadfast – the angels will descend upon them. – [Qur'an 41:30]",
			"And when you are met with a greeting, respond with a better one or return it [in a like manner]. – [Qur'an 4:86]",
			"He created the heavens and earth in truth. – [Qur'an 6:73]",
			"Say: 'He is Allah, [who is] One, Allah, the Eternal Refuge.' – [Qur'an 112:1–2]",
			"Indeed, the most noble of you in the sight of Allah is the most righteous of you. – [Qur'an 49:13]",
			"And be kind as Allah has been kind to you. – [Qur'an 28:77]",
			"Do not spy or backbite each other. – [Qur'an 49:12]",
			"Save yourselves and your families from a Fire. – [Qur'an 66:6]",
			"Paradise lies beneath the feet of mothers. – [Hadith, Nasa’i]",
			"A believer is not bitten from the same hole twice. – [Hadith, Bukhari & Muslim]",
			"Modesty brings nothing but good. – [Hadith, Bukhari & Muslim]",
			"The best jihad is speaking a word of truth to an oppressive ruler. – [Hadith, Abu Dawood]",
			"The upper hand is better than the lower hand (i.e., the giver is better than the receiver). – [Hadith, Bukhari]",
			"Beware! There is a piece of flesh in the body, and if it becomes good, the whole body becomes good. That is the heart. – [Hadith, Bukhari & Muslim]",
			"Whoever guides someone to goodness will have a reward like one who did it. – [Hadith, Muslim]",
			"Cleanliness is half of faith. – [Hadith, Muslim]",
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		quote := quotes[r.Intn(len(quotes))]

		const (
			Bold   = "\033[1m"
			Italic = "\033[3m"
			Yellow = "\033[33m"
			Reset  = "\033[0m"
		)

		fmt.Println(Bold + Yellow + "🌙 " + Italic + "❝" + quote + "❞" + Reset + " 🌙")
	}
}
