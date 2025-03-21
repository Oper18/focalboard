// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React from 'react'
import {FormattedMessage} from 'react-intl'

import './volunteersPage.scss'

const VolunteersPage = () => {
    return (
        <div className='VolunteersPage'>
            <div className='header'>
                <h2>
                    <FormattedMessage
                        id='VolunteersPage.Title'
                        defaultMessage='Volunteers'
                    />
                </h2>
            </div>
            <div className='volunteers-list'>
                {/* TODO: Add volunteers list implementation */}
            </div>
        </div>
    )
}

export default VolunteersPage
