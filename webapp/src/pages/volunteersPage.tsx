// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.
import React, {useEffect, useState} from 'react'
import {useHistory} from 'react-router-dom'

import Workspace from '../components/workspace'
import {useAppSelector} from '../store/hooks'
import {getMe} from '../store/users'
import {getVolunteers} from '../api/volunteers'
import {IVolunteer} from '../types/volunteer'

import './volunteersPage.scss'

const VolunteersContent = () => {
    const [volunteers, setVolunteers] = useState<IVolunteer[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState('')

    useEffect(() => {
        const fetchVolunteers = async () => {
            try {
                const data = await getVolunteers()
                setVolunteers(Array.isArray(data) ? data : [])
            } catch (err) {
                setError('Failed to load volunteers')
            } finally {
                setLoading(false)
            }
        }
        fetchVolunteers()
    }, [])

    if (loading) {
        return <div>{'Loading...'}</div>
    }

    if (error) {
        return <div className='error'>{error}</div>
    }

    return (
        <div className='VolunteersPage'>
            <div className='header'>
                <h2 id='volunteers-page-title'>{'Volunteers'}</h2>
            </div>
            <div className='volunteers-list'>
                {volunteers.map((volunteer) => (
                    <div
                        key={volunteer.id}
                        className='volunteer-item'
                    >
                        <h3>{volunteer.name}</h3>
                        <p>{'Contact:'} {volunteer.contact}</p>
                        {volunteer.description && <p>{volunteer.description}</p>}
                    </div>
                ))}
            </div>
        </div>
    )
}

const VolunteersPage = () => {
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
            customContent={<VolunteersContent/>}
        />
    )
}

export default VolunteersPage
