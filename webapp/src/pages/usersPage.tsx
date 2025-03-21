// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {FormattedMessage} from 'react-intl'

import './usersPage.scss'

const UsersPage = () => {
    return (
        <div className='UsersPage'>
            <div className='header'>
                <h2>
                    <FormattedMessage
                        id='UsersPage.Title'
                        defaultMessage='Users'
                    />
                </h2>
            </div>
            <div className='users-list'>
                {/* TODO: Add users list implementation */}
            </div>
        </div>
    )
}

export default UsersPage
