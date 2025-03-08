// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, {useState} from 'react'
import {Link, Redirect} from 'react-router-dom'
import {FormattedMessage} from 'react-intl'

import {useAppSelector} from '../store/hooks'
import {getLoggedIn} from '../store/users'

import './registerPage.scss'

const RegisterPage = () => {
    const errorMessage = useState('')
    const loggedIn = useAppSelector<boolean|null>(getLoggedIn)

    if (loggedIn) {
        return <Redirect to={'/'}/>
    }

    return (
        <div className='RegisterPage'>
            <Link to='/login'>
                <FormattedMessage
                    id='register.login-button-restricted'
                    defaultMessage={'Registration not available. Log in please'}
                />
            </Link>
            {errorMessage &&
                <div className='error'>
                    {errorMessage}
                </div>
            }
        </div>
    )
}

export default React.memo(RegisterPage)
