// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {FormattedMessage} from 'react-intl'
import {useHistory} from 'react-router-dom'

import Workspace from '../components/workspace'
import {useAppSelector} from '../store/hooks'
import {getMe} from '../store/users'

import './usersPage.scss'

const UsersContent = () => {
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

const UsersPage = () => {
    const history = useHistory()
    const me = useAppSelector(getMe)

    // Check if user has admin permissions
    const isAdmin = me?.permissions?.includes('manage_system')

    if (!isAdmin) {
        history.push('/')
        return null
    }

    return (
        <Workspace
            readonly={false}
            customContent={<UsersContent/>}
        />
    )
}

export default UsersPage
