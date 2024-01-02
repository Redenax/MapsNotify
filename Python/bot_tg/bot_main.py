from telegram import Update, KeyboardButton, ReplyKeyboardMarkup, ReplyKeyboardRemove
from telegram.ext import ApplicationBuilder, ContextTypes, CommandHandler, ConversationHandler, MessageHandler, filters
from connectionserver import ConnectionServer

USERNAME, PASSWORD = range(2)

NUM_TRY = 2


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE):
    keyboard = [
        [KeyboardButton("/login")],
    ]

    reply_markup = ReplyKeyboardMarkup(keyboard, resize_keyboard=True)

    chat_id = update.effective_chat.id
    await context.bot.send_message(chat_id,
                                   "Ciao " + update.message.chat.first_name + "!! Benvenuto in Traffic bot! Premi "
                                                                              "il comando di login per "
                                                                              "effettuare l'accesso.\n"
                                                                              "Digita /cancel per annullare il login.",
                                   reply_markup=reply_markup)


async def login(update: Update, context: ContextTypes.DEFAULT_TYPE):
    reply_markup = ReplyKeyboardRemove()
    await update.message.reply_text("Per favore inserisci l'email per effettuare l'accesso:\n",
                                    reply_markup=reply_markup)
    return USERNAME


async def username(update: Update, context: ContextTypes.DEFAULT_TYPE):
    email = update.message.text

    if email == '':
        await update.message.reply_text("Inserimento non valido, per favore inserisci un'email valida")
        return USERNAME

    context.user_data['email'] = email
    if 'try' not in context.user_data:
        context.user_data['try'] = 0

    await update.message.reply_text('Inserisci la password:')

    return PASSWORD


async def password(update: Update, context: ContextTypes.DEFAULT_TYPE):
    email = context.user_data['email']
    psw = update.message.text

    connect = ConnectionServer(email, psw, update.effective_chat.id)
    response = connect.connection_to_server()

    login_key = [
        [KeyboardButton("/login")],
    ]
    logout_key = [
        [KeyboardButton("/logout")],
    ]

    login_markup = ReplyKeyboardMarkup(login_key, resize_keyboard=True)

    if response == "Authorized":
        context.user_data.clear()
        context.user_data['logged'] = True
        context.user_data['connect'] = connect

        logout_markup = ReplyKeyboardMarkup(logout_key, resize_keyboard=True)

        await update.message.reply_text("Utente trovato, il servizio sta per partire")
        await update.message.reply_text("Per effettuare il logout premi logout", reply_markup=logout_markup)

        return ConversationHandler.END

    elif response == "503":
        await update.message.reply_text("Server di autenticazione momentaneamente offline.\n"
                                        "Riprovare piu tardi.", reply_markup=login_markup)
        return ConversationHandler.END

    else:
        await update.message.reply_text("Id non trovato, è possibile che tu abbia inserito un id errato o che ancora "
                                        "non esiste alcun utente con quell'ID.")

        if context.user_data['try'] != NUM_TRY:
            count = context.user_data['try']
            count += 1
            context.user_data['try'] = count

            await update.message.reply_text("Inserire nuovamente l'email: ")
            print(context.user_data['try'])

            return USERNAME

        else:
            context.user_data.clear()
            await update.message.reply_text('Numero di tentativi superato premere login per riprovare',
                                            reply_markup=login_markup)
            return ConversationHandler.END


async def handle_cancel(update: Update, context: ContextTypes.DEFAULT_TYPE):
    login_key = [
        [KeyboardButton("/login")],
    ]
    login_markup = ReplyKeyboardMarkup(login_key, resize_keyboard=True)

    await update.message.reply_text("Login interrotto!!", reply_markup=login_markup)
    return ConversationHandler.END


async def handle_logout(update: Update, context: ContextTypes.DEFAULT_TYPE):
    login_key = [
        [KeyboardButton("/login")],
    ]
    login_markup = ReplyKeyboardMarkup(login_key, resize_keyboard=True)

    if 'logged' in context.user_data and context.user_data['logged']:
        response = context.user_data['connect'].logout()
        if response == "logout effettuato":
            await update.message.reply_text("Logout effettuato!!"
                                            "Se vuoi effettuare l'accesso premi login.", reply_markup=login_markup)
            context.user_data.clear()

    else:
        await update.message.reply_text("Devi prima aver effettuato l'accesso per poter effettuare il logout.",
                                        reply_markup=login_markup)


if __name__ == '__main__':
    application = ApplicationBuilder().token('6439186304:AAE5ezRd0YgrbpCSCYSJbh_qP4DAKzlzGQ4').build()

    login_handler = ConversationHandler(
        entry_points=[CommandHandler("login", login)],
        states={
            USERNAME: [MessageHandler(filters.TEXT & ~ filters.COMMAND, username)],
            PASSWORD: [MessageHandler(filters.TEXT & ~  filters.COMMAND, password)],
        },
        fallbacks=[CommandHandler("cancel", handle_cancel)],
    )

    start_handler = CommandHandler('start', start)
    logout_handler = CommandHandler('logout', handle_logout)

    application.add_handler(start_handler)
    application.add_handler(login_handler)
    application.add_handler(logout_handler)

    application.run_polling()
